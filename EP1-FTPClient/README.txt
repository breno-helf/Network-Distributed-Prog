# EP1 Redes - Servidor FTP

## Nome: Breno Helfstein Moura
## NUSP: 9790972

## Nome: Matheus Barcellos de Castro Cunha
## NUSP: 11208238

- Compilar e executar:
    1. Para compilar os códigos fonte, o comando "make" deve ser executado dentro
    do diretório descompactado, na shell.
    2. Após compilados, será gerado o arquivo executável com o nome "ftp-server".
    3. Para executar, deve-se rodar o arquivo "ftp-server" com o argumento sendo
    uma porta para ser utilizada pelo socket "listenfd" o qual receberá conexões.
    Exemplo: "./ftp-server 8000".
    4. Para abrir conexão com o servidor, deve-se utilizar a porta passada como
    argumento no momento da execução do mesmo.

- Testes:
    Os testes rodados pela dupla foram realizados por meio do protocolo telnet e ftp 
direto na shell.

- Comandos implementados:
    USER
    PASS
    QUIT
    PWD
    CWD
    PASV
    LIST
    DELE
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

    *PWD - Utilizamos a função "getcwd()", a fim de obter o caminho do diretório em qual se 
encontra. Esse caminho e armazenado em um buffer e depois retornado ao cliente.

    *CWD - É utilizada a função "chdir()", a qual armazena em um buffer passado como argumento 
da função, o atual caminho o qual está sendo trabalhado.

    *PASV - Para gerar uma porta a fim de o cliente conectar-se ao servidor para que haja a transferên-
cia de dados, foi gerado, com auxílio da função "rand()", um número x tal que 1024<=x<=65535. Após a 
geração da porta, um socket chamado "pasvfd" e devidamente criado para poder ouvir nesta porta. Uma
mensagem é enviada para o cliente com o endereço de IP do server junto a mais duas variáveis separadas
por vírgula, no qual a multiplicação da primeira por 256, somada a segunda resulta no número da porta no
qual o socket "pasvfd" está "escutando". Exemplo de resposta ao cliente: "227 Entering Passive Mode
(127,0,0,1,251,209).".

    *LIST - Para implementar o comando, um socket chamado "datafd" e criado por meio da função "accept()"
assim que uma conexão é feita na porta passiva anteriormente criada. Após ter a conexão feita, utilizamos
a função "popen()" para dar o comando "ls -l" na shell e assim receber um ponteiro que aponta para um "FILE"
com as informações recebidas, as quais serão enviadas ao cliente usando "write()".

    *DELE - Primeiramente, é checado se o caminho fornecido como argumento existe, por meio da função
"FileOrDirExist()". Logo após a checagem de existência, outra checagem é feita a fim de saber se o
argumento fornecido e mesmo um arquivo por meio da função "FileOrDir()". Caso as duas checagem retornem
a resposta esperada, o arquivo é apagado.

    *RMD - Primeiramente, é checado se o caminho fornecido como argumento existe, por meio da função
"FileOrDirExist()". Logo após a checagem de existência, outra checagem é feita a fim de saber se o
argumento fornecido e mesmo um diretório por meio da função "FileOrDir()". Caso as duas checagem retor-
nem a resposta esperada, o diretório e apagado.

    *RETR - Faz o download do arquivo solicitado do servidor para o cliente. Para fazer o download do
arquivo, ele deve estar presente no diretório atual do servidor, caso contrário o comando retorna erro.
É implementado com uma série de comandos read e write.

    *STOR - Envia o arquivo solicitado do cliente para o servidor. Para fazer o envio do arquivo. É imple-
mentado com uma série de comandos read e write.

