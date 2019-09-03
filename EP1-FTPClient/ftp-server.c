/* Made by:
 * Breno Helfstein Moura - 9790972
 * Matheus Barcellos de Castro Cunha - 11208238
 *
 *
 * Based on example eco server provided by 
 * Prof. Daniel Batista <batista@ime.usp.br>
 *
 * TODO: don't forget to put the running instructions here.
 * 
 */

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
#include "ftp-commands.h"

int main (int argc, char **argv) {
   /* Two sockets, one that will wait for a connection and other
    * that will stabilish the connection with a specific client */
   int listenfd, connfd;
   /* The information regarding the sockets stay in this struct */
   struct sockaddr_in servaddr;

   pid_t childpid;
   /* Variables to help with the connection with the client */
   char recvline[MAXLINE + 1];
   char command[MAXDATASIZE];
   char arg[MAXDATASIZE];
   Response *res = (Response *)malloc(sizeof(Response));
   Connection *conn = (Connection *)malloc(sizeof(Connection));
   /* Store the size of the string read by the client */
   ssize_t n;

   if (argc != 2) {
      fprintf(stderr,"Use: %s <Port>\n",argv[0]);
      fprintf(stderr,"It will run a simple FTP server on the port <Port>\n");
      free(res);
      free(conn);
      exit(1);
   }

   /* Creation of a socket. It is like a file descriptor. It supports
    * operations like read, write and close. In this case the socket is 
    * created as a socket IPv4 (because of AF_INET argument) that will use
    * TCP (because of SOCK_STREAM argument) because we know that FTP runs 
    * over TCP, and will be used as a conventional application over the 
    * the internet (because of the number 0) */
   if ((listenfd = socket(AF_INET, SOCK_STREAM, 0)) == -1) {
      perror("socket :(\n");
      free(res);
      free(conn);
      exit(2);
   }

   /* Now it is necessary to inform the addresses assoctiated to this
    * socket. It is necessary to inform the address / interface and
    * the port because the socket will be waiting connection on this
    * port and address. because of this it is necessary to fill the
    * struct servaddr. It is necessary to put the type of the socket
    * (AF_INET in our case because of IPv4), in what address / interface 
    * they expect connection (Any in our cane -- INADDR_ANY) and what
    * port. In that case the port is informed as a shell argument.
    * (atoi(argv[1])). Pay attention that any port lower than 1024
    * will require sudo.
    */
   bzero(&servaddr, sizeof(servaddr));
   servaddr.sin_family      = AF_INET;
   servaddr.sin_addr.s_addr = htonl(INADDR_ANY);
   servaddr.sin_port        = htons(atoi(argv[1]));
   if (bind(listenfd, (struct sockaddr *)&servaddr, sizeof(servaddr)) == -1) {
      perror("bind :(\n");
      free(res);
      free(conn);
      exit(3);
   }

   /* As this code is a server code, the socket will be a passive socket.
    * For that it will be necessary to call the funciton listen that 
    * this socket is a server socket that will be awaiting for connection 
    * of the addressed defined by the function bind
    */
   if (listen(listenfd, LISTENQ) == -1) {
      perror("listen :(\n");
      free(res);
      free(conn);
      exit(4);
   }

   printf("[Server is online. Waiting for connection at port %s]\n",argv[1]);
   printf("[To exit, press CTRL+c or run kill/killall to kill the process]\n");
   
   /* The server is just an infinite loop waiting for connections */
   for (;;) {
      /* The initial socket that is created is the socket the will wait
       * for connection on the specified port. But it can exist several
       * clients connecting on the server. For that reason we should use
       * the accept funciton. That function will take out a connection
       * from the queue of connection that were accepted by the listenfd
       * socket and will create a specific socket for that connection.
       * The description of this new socket is returnet by the accept 
       * function. */
      if ((connfd = accept(listenfd, (struct sockaddr *) NULL, NULL)) == -1 ) {
         perror("accept :(\n");
         free(res);
         free(conn);
         exit(5);
      }
      
      /* Now the server needs to take care this client in a separate way
       * For that it is created a separate child process using the fork
       * function. The process will be a copy of this one. After the fork
       * function, both process (child and parent) will be on the same line
       * of code, but each one will have a different PID. That way it is 
       * possible to differentiate what each process will have to do. The
       * child has to processs the client request. That parent has to get
       * back on the loop to keep accepting new connections. If fork
       * returns 0 it is because it is on the child process.
       */
      if ((childpid = fork()) == 0) {
         /**** CHILD process ****/
         printf("[One connection open]\n");
         /* As we are in the child proecss we don't need listenfd
          * socket. Just parent process needs this socket. */
         close(listenfd);
         
         /* Agora pode ler do socket e escrever no socket. Isto tem
          * que ser feito em sincronia com o cliente. Não faz sentido
          * ler sem ter o que ler. Ou seja, neste caso está sendo
          * considerado que o cliente vai enviar algo para o servidor.
          * O servidor vai processar o que tiver sido enviado e vai
          * enviar uma resposta para o cliente (Que precisará estar
          * esperando por esta resposta) 
          */

         /* Now we can read from the socket and write on the socket. 
          * This has to be done in sichronization with the client. 
          * It doesn't make sense to read withut having what to read.
          * That means that in this case it is considered that the 
          * client will send something to the server. The server will
          * have to process what was sent and will send a response to
          * the client, that will have to be awaiting a response.
          */

         /* ========================================================= */
         /* ========================================================= */
         /*                         EP1 START                         */
         /* ========================================================= */
         /* ========================================================= */

         /* Making first contact and waiting for login */
         conn->socket_id = connfd;
         write(connfd, first_contact, strlen(first_contact));
         while ((n = read(connfd, recvline, MAXLINE)) > 0) {            
            recvline[n] = '\0';
            parse_ftp_command(recvline, command, arg);
            
            if (strcmp(command, "USER") != 0 && strcmp(command, "QUIT") != 0) {
               fill_message(res, "500 You must first use USER to authenticate!\n");
               client_error(connfd, res->msg);
               continue;
            }

            handle_command(command, arg, res, conn);
            if (res->error != 0) {
               client_error(connfd, res->msg);
               continue;
            }
            write_client(connfd, res->msg);

            n = read(connfd, recvline, MAXLINE);
            recvline[n] = '\0';
            parse_ftp_command(recvline, command, arg);
            if (strcmp(command, "PASS") != 0 && strcmp(command, "QUIT") != 0) {
               fill_message(res, "500 After USER command you must use PASS to authenticate!\n");
               client_error(connfd, res->msg);
               continue;               
            }
            
            handle_command(command, arg, res, conn);
            if (res->error != 0) {
               client_error(connfd, res->msg);
               continue;
            }
            write_client(connfd, res->msg);            

            break;
         }

         /* User authenticated, listening to other calls */
         while ((n=read(connfd, recvline, MAXLINE)) > 0) {
            recvline[n]='\0';
            printf("[Client connected with child process %d sent:] ",getpid());
            if ((fputs(recvline,stdout)) == EOF) {
               perror("fputs :( \n");
               free(res);
               free(conn);
               exit(6);
            }

            parse_ftp_command(recvline, command, arg);
            handle_command(command, arg, res, conn);
            if (res->error != 0) {
               client_error(connfd, res->msg);
               continue;
            }
            write_client(connfd, res->msg);            
            
         }
         /* ========================================================= */
         /* ========================================================= */
         /*                         EP1 FIM                           */
         /* ========================================================= */
         /* ========================================================= */

         /* After the connection, we can close the child process */
         printf("[One connection closed.]\n");
         free(res);
         free(conn);
         exit(0);
      }
      /**** PARENT PROCESS ****/
      /* If it is the parent process the only thing we will need to do
       * is to close the connfd, because it is the specific socket of 
       * the child process */
      close(connfd);
   }
   free(res);
   free(conn);
   exit(0);
}
