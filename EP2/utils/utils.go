package utils

/* Made by:
 * Breno Helfstein Moura - 9790972
 * Matheus Barcellos de Castro Cunha - 11208238
**/
import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

// BufferSize is the default buffer size
const BufferSize = 256

// HandlerPort is the default port for handlers
const HandlerPort = ":8042"

// GetRemoteIP extracts just the remoteIP from a connection
func GetRemoteIP(conn net.Conn) string {
	return strings.Split(conn.RemoteAddr().String(), ":")[0]
}

// UncompressChunk decompress a received chunk
func UncompressChunk(chunkStr string) (Chunk, error) {
	chunk := Chunk{}
	err := json.Unmarshal([]byte(chunkStr), &chunk)
	return chunk, err
}

// CompressChunk compress a chunk to send
func CompressChunk(chunk Chunk) (string, error) {
	chunkByte, err := json.Marshal(chunk)
	return string(chunkByte), err
}

// ClientConns fill in a channel with connections
func ClientConns(listener net.Listener) chan net.Conn {
	ch := make(chan net.Conn)
	go func() {
		for {
			client, err := listener.Accept()
			if client == nil {
				log.Printf("Couldn't accept: %v\n", err)
				continue
			}
			ch <- client
		}
	}()
	return ch
}

// GetMyIP returns your external IP and gives an error if it does not find it. -- Extracted from StackOverflow
func GetMyIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("Can't find external IP")
}

// Broadcast sends a message to all nodes. Can only be called by Master
func Broadcast(ctx *Context, msg string) error {
	if !ctx.IsMasterNode() {
		return errors.New("Only master can call Broadcast")
	}

	var wg sync.WaitGroup
	nodes := ctx.AllNodes()
	ch := make(chan bool, 5)

	for _, node := range nodes {
		ch <- true
		wg.Add(1)
		go func(ch <-chan bool, msg string, remoteIP string) {
			conn, err := net.Dial("tcp", remoteIP+HandlerPort)
			if err != nil {
				log.Println(err)
			}
			defer conn.Close()

			fmt.Fprint(conn, msg)
			<-ch
			wg.Done()
		}(ch, msg, node)
	}
	wg.Wait()

	return nil
}
