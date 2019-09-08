/*
  Utility functions for the implementation of FTP protocol
 */
#ifndef __FTP_UTILS__
#define __FTP_UTILS__

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
#define MAXDATASIZE 255
#define MAXLINE 4096
#define USERSSIZE 1

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

   /* Passive socket file descriptor, 
      -1 if we don't have any passive connection */
   int pasvfd;
} Connection;

/* Simple function to parse a FTP command line */
void parse_ftp_command(char *line, char *command, char *arg);

/* Send error message to the client and print error message on server side.
   It alsos free the memory allocate on variable msg */ 
void client_error(int connfd, char *msg);

/* Write message for client and free variable msg */
void write_client(int connfd, char *msg);

/* allocate memory for res->msg and fill it with message */
void fill_message(Response *res, const char *message);

/* Make a string all upper case */
char *turn_upper(char *str);

/* Calculate the current up address and return it with , in place of .
   Uses as base the IP of the socked in connection conn */
char *get_ip_address(Connection *conn);

/* Transform all \n of string LF_str to \r\n, it assumes that the string
   CRLF_str has enough allocated memory for that. Return the size of the 
   new string */
int transform_LF_CRLF(char *LF_str, char *CRLF_str);

#endif
