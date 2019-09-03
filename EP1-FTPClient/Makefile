# Breno Helfstein Moura       NUSP: 9790972
# Lucas Daher                 NUSP: 8991769

FLAGS = -pthread -Wall -pedantic -g -std=c99 -D _POSIX_C_SOURCE=199309L
OBJS = ftp-server.o ftp-utils.o ftp-commands.o
EXEC = ftpserver

.PHONY: clean

all: ftp-server

ftp-server: $(OBJS)
	$(CC) $(FLAGS) $^ -o $@
	make clean

%.o: %.c %.h
	$(CC) $(FLAGS) -c $< -o $@

%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

clean:
	-rm -f *.o *~
