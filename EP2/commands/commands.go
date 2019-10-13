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

	return err
}

// LEADER command will change leader
func LEADER(conn net.Conn, ctx *utils.Context, newLeader string) error {
	remoteIP := utils.GetRemoteIP(conn)
	if remoteIP != ctx.MasterNode() {
		return errors.New("Only master node can change the leader")
	}

	log.Println(fmt.Sprintf("--> Changed leader to %s", newLeader))
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
	log.Println(fmt.Sprintf("Sent chunk %d sorted!", chunk.ID))

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
		log.Println(fmt.Sprintf("--> Chunk %d sorted by %s", sortedChunk.ID, workerIP))
		_, err := fmt.Fprintf(conn, "DONE %s\n", workerIP)
		if err != nil {
			log.Printf(utils.WORKERROR, err)
		}
	case <-time.After(30 * time.Second):
		ctx.Ch() <- chunkToSort
		_, err := fmt.Fprintf(conn, "DONE %s\n", workerIP)
		if err != nil {
			log.Printf(utils.WORKERROR, err)
		}
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

	log.Println(fmt.Sprintf("--> Voting for %s!", vote))
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

	eventlog.EventFinishSorting(ctx.MasterNode())
	os.Exit(0)
	return nil
}
