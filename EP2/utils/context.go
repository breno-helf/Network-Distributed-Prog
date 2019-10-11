package utils

import "sync"

// Context is the current scenario the node is seeing the network in
type Context struct {
	nodes      map[string]bool
	leader     string
	masterNode string
	myIP       string
	ch         chan Chunk
	finalSort  bool
	mu         sync.RWMutex
}

// NewContext create a new context
func NewContext(
	nodes map[string]bool,
	leader string,
	masterNode string,
	myIP string,
	ch chan Chunk) *Context {
	return &Context{
		nodes:      nodes,
		leader:     leader,
		masterNode: masterNode,
		myIP:       myIP,
		finalSort:  false,
		ch:         ch,
	}
}

// AddNode add a new node to nodes
func (ctx *Context) AddNode(node string) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.nodes[node] = true
}

// RemoveNode remove a node from nodes
func (ctx *Context) RemoveNode(node string) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	delete(ctx.nodes, node)
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

// FinalSort returns finalSort variable
func (ctx *Context) FinalSort() bool {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	return ctx.finalSort
}

// SetFinalSort sets the final sort variable
func (ctx *Context) SetFinalSort(v bool) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.finalSort = v
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

// MyIP return myIP
func (ctx *Context) MyIP() string {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	return ctx.myIP
}
