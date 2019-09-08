# Breno Helfstein Moura						NUSP: 9790972
# Matheus Barcellos de Castro Cunha   		NUSP: 11208238

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
