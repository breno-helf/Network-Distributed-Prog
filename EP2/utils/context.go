package utils

import "sync"
import "reflect"

// Context is the current scenario the node is seeing the network in
type Context struct {
	nodes      []string
	leader     string
	masterNode string
	myIP       string
	ch         chan Chunk
	mu         sync.RWMutex
}

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

//DeleteNode Delets a node from the list
func (ctx *Context) DeleteNode(node string) {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	position := -1
	for i :=0; i < len(ctx.nodes); i++ {
		if reflect.DeepEqual(ctx.nodes[i], node){
			position = i
			break
		}
	}

	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.nodes = RemoveElement(ctx.nodes, position)

}

// CompareIP return if it's the same IP as the leader
func (ctx *Context) CompareIP(IP string) bool {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	return ctx.leader == IP
}