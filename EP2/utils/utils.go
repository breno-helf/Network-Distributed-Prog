package utils

/* Made by:
 * Breno Helfstein Moura - 9790972
 * Matheus Barcellos de Castro Cunha - 11208238
**/
import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

// BufferSize is the default buffer size
const BufferSize = 256

// HandlerPort is the default port for handlers
const HandlerPort = ":8042"

// HeartbeatTime defines Heartbeat() repeat time
const HeartbeatTime = time.Minute

// Timeout defines the timeout
const Timeout = 45 * time.Second

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
				log.Printf("Couldn't accept connection: %v\n", err)
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
				log.Printf(BROADCASTERROR, err)
				<-ch
				return
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

// CheckNode pings a node
func checkNode(conn net.Conn, remoteIP string) error {
	ch := make(chan error)
	go func(conn net.Conn) {
		_, err := fmt.Fprintf(conn, "PING\n")

		reader := bufio.NewReader(conn)
		msg, err := reader.ReadString('\n')
		if err != nil {
			ch <- err
			return
		}
		tokens := strings.Fields(msg)

		if tokens[0] != "PONG" {
			ch <- errors.New("Received message different than PONG, node is crazy (I mean, dead)")
			return
		}

		ch <- nil
	}(conn)

	select {
	case err := <-ch:
		return err
	case <-time.After(Timeout):
		return fmt.Errorf("Node %s timeouted during heartbeat", remoteIP)
	}
}

// Heartbeat Keeps HeartBeating other node
func Heartbeat(ctx *Context, remoteIP string) {
	conn, err := net.Dial("tcp", remoteIP+HandlerPort)
	if err != nil {
		log.Printf(HEARTBEATERROR, err)
		return
	}
	defer conn.Close()

	timer := time.Now()
	for {
		if time.Since(timer) >= HeartbeatTime {
			err := checkNode(conn, remoteIP)
			if err != nil {
				log.Printf(HEARTBEATERROR, err)
				Broadcast(ctx, fmt.Sprintf("DEAD %s\n", remoteIP))
			}
			timer = time.Now()
		}
	}
}

// TryEnterNetwork tries to enter the network
func TryEnterNetwork(ctx *Context) error {
	conn, err := net.Dial("tcp", ctx.MasterNode()+HandlerPort)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = fmt.Fprintf(conn, "ENTER\n")
	if err != nil {
		return err
	}

	reader := bufio.NewReader(conn)
	msg, err := reader.ReadString('\n')
	tokens := strings.Fields(msg)

	if tokens[0] == "LEADER" {
		ctx.SetLeader(tokens[1])
	} else {
		return errors.New("Expecting LEADER message")
	}

	return nil
}
