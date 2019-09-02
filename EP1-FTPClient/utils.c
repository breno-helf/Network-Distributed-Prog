#define _GNU_SOURCE
#include "utils.h"

#define LISTENQ 1
#define MAXDATASIZE 100
#define MAXLINE 4096
#define USERSSIZE 1

/* String for first contact message */
char *first_contact = "220 FTP Server (Serverzao_da_massa) [::ffff:127.0.0.1]\n";

void handle_command(char *command, char *arg, Response *res) {
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

void command_USER(char *arg, Response *res) {
   char *line;
   size_t len = 0;
   ssize_t n;
   FILE *file = fopen("./users.txt", "r");
   if (file == NULL) {
      res->username_accepted = 0;
      res->error = 1;
      sprintf(res->msg, "Could not login due to error openning users.txt file");
      return;
   }
   
   while ((n = getline(&line, &len, file)) > 0) {
      char *username;
      char *password;
      sscanf(line, "%s %s", username, password); 
      if (strcmp(username, arg) == 0) {
         
      }
   }
   
   fclose(file);
}

void fill_USER_response(Reponse *res, 

