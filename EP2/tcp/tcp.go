package tcp

/* Made by:
 * Breno Helfstein Moura - 9790972
 * Matheus Barcellos de Castro Cunha - 11208238
**/
import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"../utils"
)

const BUFFERSIZE = 1024

func ClientConns(listener net.Listener) chan net.Conn {
	ch := make(chan net.Conn)
	go func() {
		for {
			client, err := listener.Accept()
			if client == nil {
				fmt.Printf("Couldn't accept: %v\n", err)
				continue
			}
			fmt.Printf("[One connection open]\n")
			ch <- client
		}
	}()
	return ch
}

func SendFile(connection net.Conn) {
	fmt.Println("A client has connected!")
	defer connection.Close()
	file, err := os.Open("lista.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	fileSize := FillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName := FillString(fileInfo.Name(), 64)
	connection.Write([]byte(fileSize))
	connection.Write([]byte(fileName))
	sendBuffer := make([]byte, BUFFERSIZE)
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		connection.Write(sendBuffer)
	}
	return
}

func FillString(retunString string, toLength int) string {
	for {
		lengtString := len(retunString)
		if lengtString < toLength {
			retunString = retunString + ":"
			continue
		}
		break
	}
	return retunString
}

func RecieveFile() {
	connection, err := net.Dial("tcp", ":8000")
	if err != nil {
		panic(err)
	}
	defer connection.Close()
	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)

	connection.Read(bufferFileSize)
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

	connection.Read(bufferFileName)
	fileName := strings.Trim(string(bufferFileName), ":")

	newFile, err := os.Create(fileName)

	if err != nil {
		panic(err)
	}
	defer newFile.Close()
	var receivedBytes int64

	for {
		if (fileSize - receivedBytes) < BUFFERSIZE {
			io.CopyN(newFile, connection, (fileSize - receivedBytes))
			connection.Read(make([]byte, (receivedBytes+BUFFERSIZE)-fileSize))
			break
		}
		io.CopyN(newFile, connection, BUFFERSIZE)
		receivedBytes += BUFFERSIZE
	}
}

func GetLocalIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func SendChunk(connection net.Conn, r io.Reader, stop int) {
	j := 0
	sc := bufio.NewScanner(r)
	sc.Split(bufio.ScanWords)
	utils.Createfile("tmp.txt")
	fd := utils.Openfile("tmp.txt")
	var sendbuf []int
	for sc.Scan() {

		if j >= stop {
			_, erro := fd.WriteString(sc.Text())
			if erro != nil {
				fmt.Printf("Error while writing new file: %v", erro)
				return
			}
			_, erro = fd.WriteString("\n")
			if erro != nil {
				fmt.Printf("Error while writing new file: %v", erro)
				return
			}
			fd.Sync()
		} else {
			j++
			x, err := strconv.Atoi(sc.Text())
			if err != nil {
				fmt.Println("Error reading chunk file:", err)
			}
			sendbuf = append(sendbuf, x)
		}
	}

	err := os.Remove("lista.txt")
	if err != nil {
		fmt.Println("Error deleting file:", err)
	}

	err = os.Rename("tmp.txt", "lista.txt")
	if err != nil {
		fmt.Println("Error deleting file:", err)
	}

	err = binary.Write(connection, binary.BigEndian, int64(len(sendbuf)))
	if err != nil {
		fmt.Println("Error sending array lenght:", err)
	}

	for i := 0; i < len(sendbuf); i++ {
		err := binary.Write(connection, binary.BigEndian, int64(sendbuf[i]))
		if err != nil {
			fmt.Println("Error sending chunk:", err)
		}
	}
}
