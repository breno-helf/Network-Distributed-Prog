package utils

import "sync"

// Context is the current scenario the node is seeing the network in
type Context struct {
	nodes      map[string]bool
	leader     string
	masterNode string
	myIP       string
	ch         chan Chunk
	workLoad   map[string]chan bool
	mu         sync.RWMutex
	wg         sync.WaitGroup
}

// NewContext create a new context
func NewContext(
	nodes map[string]bool,
	leader string,
	masterNode string,
	myIP string,
	ch chan Chunk) *Context {
	ctx := &Context{
		nodes:      nodes,
		leader:     leader,
		masterNode: masterNode,
		myIP:       myIP,
		ch:         ch,
	}
	ctx.workLoad = make(map[string]chan bool)
	for k := range ctx.nodes {
		ctx.workLoad[k] = make(chan bool, 3)
	}

	return ctx
}

// AddNode add a new node to nodes
func (ctx *Context) AddNode(node string) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.nodes[node] = true
	ctx.workLoad[node] = make(chan bool, 3)
}

// RemoveNode remove a node from nodes
func (ctx *Context) RemoveNode(node string) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	delete(ctx.nodes, node)
	delete(ctx.workLoad, node)
}

// AllNodes returns a current snapshort from all nodes
func (ctx *Context) AllNodes() []string {
	nodes := []string{}
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	for k := range ctx.nodes {
		nodes = append(nodes, k)
	}
	return nodes
}

// WorkLoad returns the workLoad of a node
func (ctx *Context) WorkLoad(node string) (chan bool, bool) {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	ch, ok := ctx.workLoad[node]
	return ch, ok
}

// Ch return the chunk channel
func (ctx *Context) Ch() chan Chunk {
	return ctx.ch
}

// Leader returns the current leader
func (ctx *Context) Leader() string {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	return ctx.leader
}

// SetLeader changes de current leader
func (ctx *Context) SetLeader(leader string) {
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

// MyIP return myIP
func (ctx *Context) MyIP() string {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	return ctx.myIP
}

// Wg Returns the wait group
func (ctx *Context) Wg() *sync.WaitGroup {
	return &ctx.wg
}
