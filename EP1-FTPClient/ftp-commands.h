/*
 * FTP commands implementations
 */
#ifndef __FTP_COMMANDS__
#define __FTP_COMMANDS__

#define _GNU_SOURCE
#include <strings.h>
#include <arpa/inet.h>
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

/* Implements PASV command so that the client can enter passive mode
   and receive data from another address */
void command_PASV(char *arg, Response *res, Connection *conn);

/* Dummy command TYPE to respond to basic FTP client */
void command_TYPE(char *arg, Response *res, Connection *conn);

/* Implements LIST command so that the client can LIST the elements 
   in the current directory */
void command_LIST(char *arg, Response *res, Connection *conn);

/* Implements DELE command so that the client can delete a file */
void command_DELE(char *arg, Response *res, Connection *conn);

/* Implements RMD command so that the client can remove a directory */
void command_RMD(char *arg, Response *res, Connection *conn);

/* Implements RETR command so that the client can retrieve a file */
void command_RETR(char *arg, Response *res, Connection *conn);

/* Implements STOR command so that the client can store a file */
void command_STOR(char *arg, Response *res, Connection *conn);

#endif
