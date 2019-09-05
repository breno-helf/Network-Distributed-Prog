#define _GNU_SOURCE
#include "ftp-commands.h"
#include "ftp-utils.h"
#include <strings.h>
#include <string.h>


int pasvfd, connfd2;

struct sockaddr_in pasvaddr;

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

void command_PASV(char *arg, Response *res, Connection *conn) {

   if ((pasvfd = socket(AF_INET, SOCK_STREAM, 0)) == -1) {
      perror("socket :(\n");
      exit(2);
   }

   srand (time(NULL));
   int connect_port = ((rand()%64511)+1024);
   bzero(&pasvaddr, sizeof(pasvaddr));
   pasvaddr.sin_family      = AF_INET;
   pasvaddr.sin_addr.s_addr = htonl(INADDR_ANY);
   pasvaddr.sin_port        = htons(connect_port);

   if (bind(pasvfd, (struct sockaddr *)&pasvaddr, sizeof(pasvaddr)) == -1) {
      perror("bind :(\n");
      exit(3);
   }

   printf("%d\n",connect_port);

   if (listen(pasvfd, LISTENQ) == -1) {
      perror("listen :(\n");
      exit(4);
   }

   if ((connfd2 = accept(pasvfd, (struct sockaddr *) NULL, NULL)) == -1 ) {
      perror("accept :(\n");
      exit(5);
   }

   printf("yayayay\n");
}

void command_TYPE(char *arg, Response *res, Connection *conn) {
   res->error = 0;
   fill_message(res, "200 Type is just a dummy command for this recreational FTP\n");
}
