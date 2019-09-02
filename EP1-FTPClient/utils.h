#define _GNU_SOURCE
#include <stdio.h>
#include <stdlib.h>
#include <errno.h>
#include <string.h>
#include <netdb.h>
#include <sys/types.h>
#include <netinet/in.h>
#include <sys/socket.h>
#include <arpa/inet.h>
#include <time.h>
#include <unistd.h>

#define LISTENQ 1
#define MAXDATASIZE 100
#define MAXLINE 4096
#define USERSSIZE 1

typedef struct Response {
   /* Response message */
   char *msg = NULL;

   /* 1 if the command had an error, 0 otherwise */
   int error = 0;
} Response;

typedef struct Connection {
   char *username;

   /* Socket of the current connection */
   int socket_id;   
} Connection;

/* String for first contact message */
char *first_contact = "220 FTP Server (Serverzao_da_massa) [::ffff:127.0.0.1]\n";

void handle_command(char *command, char *arg, Response *res, Connection *conn);

void command_USER(char *arg, Response *res, Connection *conn);

void command_PASS(char *arg, Response *res, Connection *conn);

void command_QUIT(char *arg, Response *res, Connection *conn);
