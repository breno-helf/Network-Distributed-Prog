# EP1 Redes - Servidor FTP

## Nome: Breno Helfstein Moura
## NUSP: 9790972

## Nome: Matheus Barcellos de Castro Cunha
## NUSP: 11208238

- Comandos implementados:
    USER
    PASS
    QUIT
    PWD
    CWD
    PASV
    LIST DELE
    RMD
    RETR
    STOR

- Login no "Serverzao_da_massa":
    Para autenticar no servidor, basta escolher qualquer combinação de
usuario e senha. Por exemplo: 'USER a', 'PASS a'.

- Resumo das implementações dos comandos:
    *USER - É feita uma checagem a fim de conferir se o argumento passado como 
usuário tem o valor igual a NULL, e também é checado se o cliente do processo x 
já havia realizado login.

    *PASS - É feita uma checagem a fim de conferir se o argumento passado como 
senha tem o valor igual a NULL, e também é checado se o cliente do processo x 
já havia realizado login.

    *QUIT - É enviada uma mensagem de fim de conexão para o cliente pelo socket "connfd".
Todos os recursos associados ao socket "connfd" são liberados e sua conexão é fechada, idem
a eventuais sockets criados para conexões em modo passivo.

