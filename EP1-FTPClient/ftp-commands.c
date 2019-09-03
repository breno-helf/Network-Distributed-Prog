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
   } else {
      res->error = 1;
      res->msg = malloc(sizeof(char) * MAXDATASIZE);
      sprintf(res->msg, "500 %s not understood\n", command);
   }

   /*
     TO IMPLEMENT, PLEASE ADD THE COMMAND THAT WE NEED TO IMPLEMENT HERE.     
     
   else if (strcmp(command, "LIST") == 0) {
      NULL;
   } else if (strcmp(command, "DELE") == 0) {
      NULL;
   } else if (strcmp(command, "STOR") == 0) {
      NULL;
   } else if (strcmp(command, "PASV") == 0) {
      NULL;
   } else if (strcmp(command, "RETR") == 0) {
      NULL;
   }

   */
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
   conn->username = malloc(sizeof(char) * strlen(arg));
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
   fill_message(res, "221 Goodbye\n");
   write(conn->socket_id, res->msg, strlen(res->msg));
   fprintf(stderr, "[Client %d] - %s\n", conn->socket_id, res->msg);
   close(conn->socket_id);
   free(res->msg);
   free(conn->username);
   exit(0);
}

void command_PWD(char *arg, Response *res, Connection *conn) {
   char path_name[256];
   getcwd(path_name, sizeof(path_name));
   sprintf(res->msg, "257 \"%s\" is the curent directory\n", path_name);
}

void command_CWD(char *arg, Response *res, Connection *conn) {
   if (chdir(arg) == 0) {
      sprintf(res->msg, "250 CWD command successful\n");
   }
   else {
      sprintf(res->msg, "550 \"%s\": No such file or directory\n", arg);
   }
}

char *turn_upper(char *str) {
  unsigned char *p = (unsigned char *)str;
  while (*p) {
     *p = toupper(*p);
      p++;
  }
  return str;
}
