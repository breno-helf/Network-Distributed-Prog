/*
 * FTP commands implementations
 */
#ifndef __FTP_COMMANDS__
#define __FTP_COMMANDS__

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
#include "ftp-utils.h"

/* Gets a parsed command and executes the specified action, filling response
   variable res. */
void handle_command(char *command, char *arg, Response *res, Connection *conn);

/* Implements authentication command USER, for the client to specify what
   user he will use to log in */
void command_USER(char *arg, Response *res, Connection *conn);

/* Implements authentication command PASS, for the client to specify what
   is the password for the user he speificied */
void command_PASS(char *arg, Response *res, Connection *conn);

/* Implements QUIT command so that the client can close the connection */
void command_QUIT(char *arg, Response *res, Connection *conn);

/* Implements PWD command so that the client can check which path he is on */
void command_PWD(char *arg, Response *res, Connection *conn);

/* Implements PWD command so that the client can change his path */
void command_CWD(char *arg, Response *res, Connection *conn);

#endif
