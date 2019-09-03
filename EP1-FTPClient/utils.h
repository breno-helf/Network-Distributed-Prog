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
#include <ctype.h>

#define LISTENQ 1
#define MAXDATASIZE 100
#define MAXLINE 4096
#define USERSSIZE 1
#define MAXSTRINGSIZE 280

typedef struct Response {
   /* Response message */
   char *msg;

   /* 1 if the command had an error, 0 otherwise */
   int error;
} Response;

typedef struct Connection {
   /* Username logged in */
   char *username;

   /* 1 if the user is logged in, 0 otherwise */
   int logged_status;
   
   /* Socket of the current connection */
   int socket_id;   
} Connection;

char path_name[256];

/* String for first contact message */
static char *first_contact = "220 FTP Server (Serverzao_da_massa) [::ffff:127.0.0.1]\n";

void handle_command(char *command, char *arg, Response *res, Connection *conn);

void command_USER(char *arg, Response *res, Connection *conn);

void command_PASS(char *arg, Response *res, Connection *conn);

void command_QUIT(char *arg, Response *res, Connection *conn);

void command_PWD(char *arg, Response *res, Connection *conn);

void command_CWD(char *arg, Response *res, Connection *conn);

char *turn_upper(char *str);
