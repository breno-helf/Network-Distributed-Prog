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
}

void write_client(int connfd, char *msg) {
   write(connfd, msg, strlen(msg));
}

void fill_message(Response *res, const char *message) {
   res->msg = (char *)malloc(sizeof(char) * strlen(message) + 1);
   strcpy(res->msg, message);
}

char *turn_upper(char *str) {
   unsigned char *p = (unsigned char *)str;
   while (*p) {
      *p = toupper(*p);
      p++;
   }
   return str;
}

char *get_ip_address(Connection *conn) {
   socklen_t addr_size = sizeof(struct sockaddr_in);
   struct sockaddr_in addr;
   getsockname(conn->socket_id, (struct sockaddr *)&addr, &addr_size);
 
   return inet_ntoa(addr.sin_addr);;    
}

int transform_LF_CRLF(char *LF_str, char *CRLF_str) {
   int i, j;
   for (i = 0, j = 0; LF_str[i] != '\0'; i++, j++) {
      if (LF_str[i] == '\n') {
         CRLF_str[j++] = '\r';
      }
      CRLF_str[j] = LF_str[i];      
   }
   CRLF_str[j] = '\0';
   
   return j;
}

int FileOrDir(const char* FileDir) {
   struct stat path;
   stat(FileDir, &path);
   return S_ISDIR(path.st_mode);
}

int FileOrDirExist (const char* FileDir) {
   struct stat path;
   return lstat(FileDir,&path);
}
