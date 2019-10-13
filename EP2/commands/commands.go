package commands

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"sort"
	"strings"
	"time"

	"../eventlog"
	"../leader"
	"../utils"
)

// ENTER command will allow someone to enter in the network
func ENTER(conn net.Conn, ctx *utils.Context) error {
	if !ctx.IsMasterNode() {
		return errors.New("Can't let someone enter if I am not the master node")
	}

	remoteIP := utils.GetRemoteIP(conn)
	fmt.Println("Remote address entering network", remoteIP)

	_, err := fmt.Fprintf(conn, "LEADER %s\n", ctx.Leader())
	if err != nil {
		return err
	}

	nodesSlice := ctx.AllNodes()
	nodesBytes, err := json.Marshal(nodesSlice)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(conn, "NODES %s\n", string(nodesBytes))
	if err != nil {
		return err
	}

	err = utils.Broadcast(ctx, fmt.Sprintf("ENTERED %s\n", remoteIP))

	go utils.Heartbeat(ctx, remoteIP)

	return err
}

// LEADER command will change leader
func LEADER(conn net.Conn, ctx *utils.Context, newLeader string) error {
	remoteIP := utils.GetRemoteIP(conn)
	if remoteIP != ctx.MasterNode() {
		return errors.New("Only master node can change the leader")
	}

	ctx.SetLeader(newLeader)
	eventlog.EventLeaderElected(newLeader)
	if newLeader == ctx.MyIP() {
		go leader.Leader(ctx)
	}

	return nil
}

// SORT received a chunk, and decompress it sorting and sent it back to the master
func SORT(conn net.Conn, ctx *utils.Context, chunkStr string) error {
	remoteIP := utils.GetRemoteIP(conn)
	if remoteIP != ctx.MasterNode() {
		return errors.New("Only master node can send an array for sorting")
	}

	chunk, err := utils.UncompressChunk(chunkStr)
	if err != nil {
		return err
	}
	sort.Ints(chunk.Numbers)

	chunkStr, err = utils.CompressChunk(chunk)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(conn, "SORTED %s\n", chunkStr)

	return err
}

// WORK will receive an IP that is requesting work.
// If master will send an array for sorting
func WORK(conn net.Conn, ctx *utils.Context, workerIP string) error {
	if !ctx.IsMasterNode() {
		return errors.New("Only master node can receive a WORK order")
	}

	remoteIP := utils.GetRemoteIP(conn)
	if remoteIP != ctx.Leader() {
		return errors.New("Only leader can send a WORK order")
	}

	chunkToSort, ok := <-ctx.Ch()
	if !ok {
		// There is no chunk to sort
		return nil
	}

	ch := make(chan utils.Chunk, 1)
	go func(ch chan utils.Chunk, workerIP string) {
		workerConn, err := net.Dial("tcp", workerIP+utils.HandlerPort)
		if err != nil {
			log.Printf(utils.WORKERROR, err)
			return
		}
		defer workerConn.Close()

		chunkStr, err := utils.CompressChunk(chunkToSort)
		if err != nil {
			log.Printf(utils.WORKERROR, err)
			return
		}

		_, err = fmt.Fprintf(workerConn, "SORT %s\n", chunkStr)
		if err != nil {
			log.Printf(utils.WORKERROR, err)
			return
		}
		reader := bufio.NewReader(workerConn)
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Printf(utils.WORKERROR, err)
			return
		}
		tokens := strings.Fields(msg)

		if tokens[0] != "SORTED" {
			log.Printf(utils.WORKERROR, errors.New("Received message different than SORTED"))
			return
		}

		sortedChunk, err := utils.UncompressChunk(tokens[1])

		ch <- sortedChunk
	}(ch, workerIP)

	select {
	case sortedChunk := <-ch:
		utils.StoreChunk(sortedChunk)
		ctx.Wg().Done()
		fmt.Fprintf(conn, "DONE %s\n", workerIP)
		log.Printf("Machine %s ordered chunk %d.", workerIP, sortedChunk.ID)
	case <-time.After(utils.Timeout):
		ctx.Ch() <- chunkToSort
		fmt.Fprintf(conn, "DONE %s\n", workerIP)
		return fmt.Errorf("TIMEOUT: Machine %s timeouted during sorting of chunk %d", workerIP, chunkToSort.ID)
	}

	return nil
}

// ENTERED acknowledges the entrance of a new node
func ENTERED(conn net.Conn, ctx *utils.Context, newNode string) error {
	remoteIP := utils.GetRemoteIP(conn)
	if remoteIP != ctx.MasterNode() {
		return errors.New("Only master can communicate the entrance of a new node")
	}

	ctx.AddNode(newNode)
	eventlog.EventNewNode(newNode)

	return nil
}

// ELECTION is the calling of an election
func ELECTION(conn net.Conn, ctx *utils.Context) error {
	remoteIP := utils.GetRemoteIP(conn)
	if remoteIP != ctx.MasterNode() {
		return errors.New("Only master can call election")
	}

	eventlog.EventElectingLeader()
	ctx.SetLeader("")
	nodes := ctx.AllNodes()
	vote := nodes[rand.Intn(len(nodes))]

	_, err := fmt.Fprintf(conn, "VOTE %s\n", vote)

	return err
}

// NODES inform all the nodes at once
func NODES(conn net.Conn, ctx *utils.Context, nodesStr string) error {
	remoteIP := utils.GetRemoteIP(conn)
	if remoteIP != ctx.MasterNode() {
		return errors.New("Only master can update nodes")
	}

	nodes := []string{}
	err := json.Unmarshal([]byte(nodesStr), &nodes)
	if err != nil {
		return err
	}

	for _, node := range nodes {
		ctx.AddNode(node)
	}

	return nil
}

// END finishes the sorting
func END(conn net.Conn, ctx *utils.Context) error {
	remoteIP := utils.GetRemoteIP(conn)
	if remoteIP != ctx.MasterNode() {
		return errors.New("Only master can end sorting")
	}

	if ctx.IsMasterNode() {
		return nil // MasterNode kill itself
	}

	eventlog.EventFinishSorting(ctx.MasterNode())
	os.Exit(0)
	return nil
}

// DEAD recognizes that a node is dead
func DEAD(conn net.Conn, ctx *utils.Context, deadNode string) error {
	remoteIP := utils.GetRemoteIP(conn)
	if remoteIP != ctx.MasterNode() {
		return errors.New("Only master can inform a dead node")
	}

	var err error
	eventlog.EventDeadNode(deadNode)
	ctx.RemoveNode(deadNode)
	if deadNode == ctx.Leader() && ctx.IsMasterNode() {
		deadLeaderCh := ctx.DeadLeaderCh()
		select {
		case deadLeaderCh <- true:
		default:
			err = errors.New("DeadLeader Channel is full")
		}
	}

	if deadNode == ctx.MyIP() {
		go enterNetwork(ctx)
	}

	return err
}

// EnterNetwork enter the network
func enterNetwork(ctx *utils.Context) {
	err := utils.TryEnterNetwork(ctx)
	for err != nil {
		err = utils.TryEnterNetwork(ctx)
	}
	ctx.AddNode(ctx.MyIP())
	eventlog.EventNewNode(ctx.MyIP())
}

// PING returns PONG
func PING(conn net.Conn, ctx *utils.Context) error {
	_, err := fmt.Fprintf(conn, "PONG\n")
	return err
}
