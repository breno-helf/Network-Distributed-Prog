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

// Chunk is the chunk to be ordered by some node
type Chunk struct {
	Numbers []int
	ID      int
}

// Context is the current scenario the node is seeing the network in
type Context struct {
	nodes      []string
	leader     string
	masterNode string
	myIP       string
	ch         chan Chunk
	mu         sync.RWMutex
}

// BufferSize is the default buffer size
const BufferSize = 256

// HandlerPort is the default port for handlers
const HandlerPort = ":8042"

// NewContext create a new context
func NewContext(
	nodes []string,
	leader string,
	masterNode string,
	myIP string,
	ch chan Chunk) *Context {
	return &Context{
		nodes:      nodes,
		leader:     leader,
		masterNode: masterNode,
		myIP:       myIP,
		ch:         ch,
	}
}

// AddNode append a new node to the Nodes slice
func (ctx *Context) AddNode(node string) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.nodes = append(ctx.nodes, node)
}

// Leader returns the current leader
func (ctx *Context) Leader() string {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	return ctx.leader
}

// ChangeLeader changes de current leader
func (ctx *Context) ChangeLeader(leader string) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.leader = leader
}

// MasterNode returns the master node
func (ctx *Context) MasterNode() string {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	return ctx.masterNode
}

// IsMasterNode return if it is a master node
func (ctx *Context) IsMasterNode() bool {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	return ctx.masterNode == ctx.myIP
}

// IsLeader return if it is the leader
func (ctx *Context) IsLeader() bool {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	return ctx.leader == ctx.myIP
}

// GetRemoteIP extracts just the remoteIP from a connection
func GetRemoteIP(conn net.Conn) string {
	return strings.Split(conn.RemoteAddr().String(), ":")[0]
}

// UncompressChunk decompress a received chunk
func UncompressChunk(chunk string) ([]int, error) {
	chunkSlice := []int{}
	err := json.Unmarshal([]byte(chunk), &chunkSlice)
	return chunkSlice, err
}

// CompressChunk compress a chunk to send
func CompressChunk(chunkSlice []int) string {
	chunk, _ := json.Marshal(chunkSlice)
	return string(chunk)
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
			fmt.Printf("[One connection open]\n")
			ch <- client
		}
	}()
	return ch
}

// GetMyIP returns your external IP and gives an error if it does not find it.
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
