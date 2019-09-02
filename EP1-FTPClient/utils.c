#define _GNU_SOURCE
#include "utils.h"

#define LISTENQ 1
#define MAXDATASIZE 100
#define MAXLINE 4096
#define USERSSIZE 1

void handle_command(char *command, char *arg, Response *res, Connection *conn) {
   if (strcmp(command, "USER") == 0) {
      command_USER(arg, res);
   } else if (strcmp(command, "PASS") == 0) {
      command_PASS(arg, res);
   } else if (strcmp(command, "LIST") == 0) {
      NULL;
   } else if (strcmp(command, "DELE") == 0) {
      NULL;
   } else if (strcmp(command, "STOR") == 0) {
      NULL;
   } else if (strcmp(command, "QUIT") == 0) {
      command_QUIT(arg, res);
   }
}

void command_USER(char *arg, Response *res, Connection *conn) {
   if (arg == NULL) {
      sprintf(res->msg, "500 USER: command requires a parameter\n");
      res->error = 0;
      return;
   }
   
   sprintf(res->msg, "331 Password required for %s\n", arg);
   res->error = 0;
}
