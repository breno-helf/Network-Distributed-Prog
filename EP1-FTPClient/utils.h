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
   /* Name of the command answering to */
   char *cmd;

   /* 1 if username is accepted, 0 otherwise */
   int username_accepted;

   /* Username logged in or currently trying to log in */ 
   char *username;
   
   /* 1 if logged in, 0 otherwise */
   int logged_in;

   /* Response message */
   char *msg;

   /* 1 if the command had an error, 0 otherwise */
   int error;
} Response;

/* String for first contact message */
char *first_contact = "220 FTP Server (Serverzao_da_massa) [::ffff:127.0.0.1]\n";

void handle_command(char *command, char *arg, Response *res);

void command_USER(char *arg, Response *res);

void command_PASS(char *arg, Response *res);

void command_QUIT(char *arg, Response *res);
