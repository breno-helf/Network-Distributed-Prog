#define _GNU_SOURCE
#include "ftp-commands.h"

void handle_command(char *command, char *arg, Response *res, Connection *conn) {
   if (strcmp(command, "USER") == 0) {
      command_USER(arg, res, conn);
   } else if (strcmp(command, "PASS") == 0) {
      command_PASS(arg, res, conn);
   } else if (strcmp(command, "QUIT") == 0) {
      command_QUIT(arg, res, conn);
   } else if (strcmp(command, "PWD") == 0) {
      command_PWD(arg, res, conn);
   } else if (strcmp(command, "CWD") == 0) {
      command_CWD(arg, res, conn);
   } else if (strcmp(command, "PASV") == 0) {
      command_PASV(arg, res, conn);
   } else if (strcmp(command, "LIST") == 0) {
      command_LIST(arg, res, conn);
   } else if (strcmp(command, "DELE") == 0) {
      command_DELE(arg, res, conn);
   } else if (strcmp(command, "RMD") == 0)  {
      command_RMD(arg, res, conn);
   } else if (strcmp(command, "RETR") == 0) {
      command_RETR(arg, res, conn);
   } else if (strcmp(command, "STOR") == 0) {
      command_STOR(arg, res, conn);
   } else {
      res->error = 1;
      res->msg = malloc(sizeof(char) * MAXDATASIZE);
      sprintf(res->msg, "500 %s not understood\n", command);
   }
}

void command_USER(char *arg, Response *res, Connection *conn) {   
   if (arg == NULL) {      
      fill_message(res, "500 USER: command requires a parameter\n");
      res->error = 1;
      return;
   }

   if (conn->logged_status == 1) {
      fill_message(res, "501 Reauthentication not supported\n");
      res->error = 1;
      return;
   }

   res->msg = malloc(sizeof(char) * MAXDATASIZE);
   sprintf(res->msg, "331 Password required for %s\n", arg);
   res->error = 0;
   conn->username = malloc(sizeof(char) * MAXDATASIZE);
   strcpy(conn->username, arg);
}

void command_PASS(char *arg, Response *res, Connection *conn) {
   if (arg == NULL) {
      fill_message(res, "500 PASS: command requires a parameter\n");
      res->error = 1;
      return;   
   }

   if (conn->logged_status == 1) {
      fill_message(res, "503 You are already logged in\n");
      res->error = 1;
      return;
   }
   
   res->msg = malloc(sizeof(char) * MAXDATASIZE);
   sprintf(res->msg, "230 User %s logged in\n", conn->username);
   res->error = 0;
   conn->logged_status = 1;
}

void command_QUIT(char *arg, Response *res, Connection *conn) {
   res->error = 0;
   fill_message(res, "221 Goodbye\n");
   write(conn->socket_id, res->msg, strlen(res->msg));
   fprintf(stderr, "[Client %d] - %s\n", conn->socket_id, res->msg);
   close(conn->socket_id);
   if (conn->pasvfd >= 0) {
      close(conn->pasvfd);
      conn->pasvfd = -1;
   }
   free(res->msg);
   free(conn->username);
   exit(0);
}

void command_PWD(char *arg, Response *res, Connection *conn) {
   char path_name[MAXDATASIZE];

   if (getcwd(path_name, sizeof(path_name)) == NULL) {
      res->error = 1;
      fill_message(res, "500 Failed to get the path");
      return;
   } 

   res->error = 0;
   res->msg = malloc(sizeof(char) * MAXDATASIZE);
   sprintf(res->msg, "257 \"%s\" is the curent directory\n", path_name);
}

void command_CWD(char *arg, Response *res, Connection *conn) {
   if (chdir(arg) == 0) {
      res->error = 0;
      fill_message(res, "250 CWD command successful\n");
      return;
   }
   
   res->error = 1;
   res->msg = malloc(sizeof(char) * MAXDATASIZE);
   sprintf(res->msg, "550 \"%s\": No such file or directory\n", arg);
}

void command_PASV(char *arg, Response *res, Connection *conn) {
   if ((conn->pasvfd = socket(AF_INET, SOCK_STREAM, 0)) == -1) {
      res->error = 1;
      fill_message(res, "500 Failed to create a socket\n");
      return;
   }

   srand(time(NULL));
   int connect_port = (rand() % 64511) + 1024;
   struct sockaddr_in pasvaddr;
   bzero(&pasvaddr, sizeof(pasvaddr));
   pasvaddr.sin_family = AF_INET;
   pasvaddr.sin_addr.s_addr = htonl(INADDR_ANY);
   pasvaddr.sin_port = htons(connect_port);

   if (bind(conn->pasvfd, (struct sockaddr *)&pasvaddr, sizeof(pasvaddr)) == -1) {
      res->error = 1;
      fill_message(res, "500 Failed to bind socket\n");
      return;
   }

   if (listen(conn->pasvfd, LISTENQ) == -1) {
      res->error = 1;
      fill_message(res, "500 Failed make socket listen\n");
      return;      
   }
   
   res->error = 0;
   res->msg = malloc(sizeof(char) * MAXDATASIZE); 
   char *ip_address = get_ip_adddress(conn);

   sprintf(res->msg,
           "227 Entering Passive Mode (%s,%d,%d).\n",
           ip_address, connect_port / 256, connect_port % 256);
}

void command_TYPE(char *arg, Response *res, Connection *conn) {
   res->error = 0;
   fill_message(res, "200 Type is just a dummy command for this recreational FTP\n");
}

void command_LIST(char *arg, Response *res, Connection *conn) {
   int datafd;
   
   if (conn->pasvfd != -1) {
      /* We are in passive mode */
      if ((datafd = accept(conn->pasvfd, (struct sockaddr *)NULL, NULL)) == -1) {
         res->error = 1;
         fill_message(res, "500 Failed to connect to server\n");
         return;
      }
   } else {
      res->error = 1;
      fill_message(res, "500 Must be in passive mode\n");
      return;
   }

   write_client(conn->socket_id, "150 Opening ASCII mode data connection for file list\n");
    
   char path_name[MAXDATASIZE];
   char buffer[MAXDATASIZE];
   char file_buffer[MAXDATASIZE];
   char CRLF_file_buffer[2 * MAXDATASIZE];
   /* How much we've read */
   int n;
   int first_line = -1;
   int bytes_sent;
   
   res->error = 0;
   getcwd(path_name, sizeof(path_name));
   sprintf(buffer, "ls -l %s", path_name);
   FILE *p = popen(buffer, "r");
   while ((n = fread(file_buffer, 1, MAXDATASIZE, p)) > 0) {
      file_buffer[n] = '\0';
      
      int m = transform_LF_CRLF(file_buffer, CRLF_file_buffer);
      if (first_line == -1) {
         char *newline = strchr(CRLF_file_buffer, '\n');
         if (newline != NULL) {
            first_line = newline - CRLF_file_buffer;
            bytes_sent = write(datafd, newline + 1, m - first_line - 1);
         } else {
            continue;
         }
      } else {
         bytes_sent = write(datafd, CRLF_file_buffer, m);
      }
      
      if (bytes_sent == -1) {
         res->error = 1;
         fill_message(res, "500 We failed to sent bytes to the client side\n");
         break;
      }
   }
   
   if (res->error == 0)
      fill_message(res, "226 Transfer complete\n");
   
   pclose(p);
   close(datafd);
   if (conn->pasvfd >= 0) {
      close(conn->pasvfd);
      close(datafd);
      conn->pasvfd = -1;
   }
}

void command_DELE(char *arg, Response *res, Connection *conn) {

   if (FileOrDirExist(arg) != 0) {
      res->error = 1;
      char msg[MAXDATASIZE];
      sprintf(msg, "550 %s: No such file or directory\n", arg);
      fill_message(res, msg);
      return;
   }

   if (FileOrDir(arg) == 1) {
         res->error = 1;
         char msg[MAXDATASIZE];
         sprintf(msg, "550 %s: Not a file\n", arg);
         fill_message(res, msg);
         return;
   }

   remove(arg);
   res->error = 0;
   fill_message(res, "250 DELE command successful\n");
   return;
}

void command_RMD(char *arg, Response *res, Connection *conn) {

   if (FileOrDirExist(arg) != 0) {
      res->error = 1;
      char msg[MAXDATASIZE];
      sprintf(msg, "550 %s: No such file or directory\n", arg);
      fill_message(res, msg);
      return;
   }

   if (FileOrDir(arg) == 0) {
      res->error = 1;
      char msg[MAXDATASIZE];
      sprintf(msg, "550 %s: Not a directory\n", arg);
      fill_message(res, msg);
      return;
   }

   remove(arg);
   res->error = 0;
   fill_message(res, "250 RMD command successful\n");
   return;
}

void command_RETR(char *arg, Response *res, Connection *conn) {   
   int datafd;

   if (conn->pasvfd != -1) {
      /* We are in passive mode */
      if ((datafd = accept(conn->pasvfd, (struct sockaddr *)NULL, NULL)) == -1) {
         res->error = 1;
         fill_message(res, "500 Failed to connect to server\n");
         return;
      }
   } else {
      res->error = 1;
      fill_message(res, "500 Must be in passive mode\n");
      return;
   }

   char buffer[1024];
   /* How much we've read */
   int n;

   res->error = 0;
   FILE *file = fopen(arg, "r");

   if (file == NULL) {
      char error_msg [MAXDATASIZE];
      res->error = 1;
      sprintf(error_msg,"550 %s: No such file or directory\n", arg);
      fill_message(res, error_msg);
      close(datafd);
      if (conn->pasvfd >= 0) {
         close(conn->pasvfd);
         conn->pasvfd = -1;
      }
      return;
   }

   char msg[MAXDATASIZE];

   sprintf(msg, "150 Opening BINARY mode data connection for %s\n", arg);

   write_client(conn->socket_id, msg);

   int filefd = fileno(file);

   while ((n = read(filefd, buffer, sizeof(buffer))) > 0) {
      int bytes_written = write(datafd, buffer, n); 
      if (bytes_written == -1) {
         res->error = 1;
         fill_message(res, "500 We failed to write bytes to the socket\n");
         break;
      }
   }

   if (res->error == 0) {
      fill_message(res, "226 Transfer complete\n");
   }

   fclose(file);
   if (conn->pasvfd >= 0) {
      close(conn->pasvfd);
      close(datafd);
      conn->pasvfd = -1;
   }
}

void command_STOR(char *arg, Response *res, Connection *conn) {
   int datafd;
   
   if (conn->pasvfd != -1) {
      /* We are in passive mode */
      if ((datafd = accept(conn->pasvfd, (struct sockaddr *)NULL, NULL)) == -1) {
         res->error = 1;
         fill_message(res, "500 Failed to connect to server\n");
         return;
      }
   } else {
      res->error = 1;
      fill_message(res, "500 Must be in passive mode\n");
      return;
   }

   char buffer[1024];
   /* How much we've read */
   int n;

   res->error = 0;
   FILE *file = fopen(arg, "w");

   if (file == NULL) {
      res->error = 1;
      fill_message(res, "451 Could not open file\n");
      close(datafd);
      if (conn->pasvfd >= 0) {
         close(conn->pasvfd);
         conn->pasvfd = -1;
      }
      return;
   }

   char msg[MAXDATASIZE];

   sprintf(msg, "150 Opening BINARY mode data connection for %s\n", arg);

   write_client(conn->socket_id, msg);

   int filefd = fileno(file);
   while ((n = read(datafd, buffer, sizeof(buffer))) > 0) {
      int bytes_written = write(filefd, buffer, n); 
      if (bytes_written == -1) {
         res->error = 1;
         fill_message(res, "500 We failed to write bytes to the file\n");
         break;
      }
   }

   if (res->error == 0)
      fill_message(res, "226 Transfer complete\n");
   
   fclose(file);
   close(datafd);
   if (conn->pasvfd >= 0) {
      close(conn->pasvfd);
      conn->pasvfd = -1;
   }
}
