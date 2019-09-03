#define _GNU_SOURCE
#include "ftp-utils.h"

void parse_ftp_command(char *line, char *command, char *arg) {
   /* In case we have an empty command, we shall have empty command and arg */
   strcpy(command, "");
   strcpy(arg, "");
   sscanf(line, "%s %s", command, arg);
}

void client_error(int connfd, char *msg) {
   write(connfd, msg, strlen(msg));
   fprintf(stderr, "[Client %d] - ERROR: %s\n", connfd, msg);
   free(msg);
}

void write_client(int connfd, char *msg) {
   write(connfd, msg, strlen(msg));
   free(msg);
}

void fill_message(Response *res, const char *message) {
   res->msg = (char *)malloc(sizeof(char) * strlen(message) + 1);
   strcpy(res->msg, message);
}
