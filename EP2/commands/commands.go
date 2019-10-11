package commands

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"sort"
	"strings"
	"time"

	"../eventlog"
	"../utils"
)

// ENTER command will allow someone to enter in the network
func ENTER(conn net.Conn, ctx *utils.Context) error {
	if !ctx.IsMasterNode() {
		return errors.New("Can't let someone enter if I am not the master node")
	}

	remoteIP := utils.GetRemoteIP(conn)
	fmt.Println("Remote address entering network", remoteIP)
	ctx.AddNode(remoteIP)

	_, err := fmt.Fprintf(conn, "LEADER %s\n", ctx.Leader())

	if err != nil {
		return err
	}

	eventlog.EventNewNode(remoteIP)

	return nil
}

// LEADER command will change leader
func LEADER(conn net.Conn, ctx *utils.Context, newLeader string) error {
	remoteIP := utils.GetRemoteIP(conn)
	if remoteIP != ctx.MasterNode() {
		return errors.New("Only master node can change the leader")
	}

	ctx.ChangeLeader(newLeader)
	eventlog.EventLeaderElected(newLeader)

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
	if err != nil {
		return err
	}

	return nil
}

// WORK will receive an IP that is requesting work.
// If master will send an array for sorting
func WORK(conn net.Conn, ctx *utils.Context, remoteIP string) error {
	if !ctx.IsMasterNode() {
		return errors.New("Only master node can receive a WORK order")
	}

	chunkToSort, ok := <-ctx.Ch()
	if !ok {
		ctx.SetFinalSort(true)
		// There is no chunk to sort
		return nil
	}

	ch := make(chan utils.Chunk, 1)
	go func(ch chan utils.Chunk, remoteIP string) {
		workerConn, err := net.Dial("tcp", remoteIP+utils.HandlerPort)
		if err != nil {
			log.Println(err)
			return
		}

		chunkStr, err := utils.CompressChunk(chunkToSort)
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Fprintf(workerConn, "SORT %s\n", chunkStr)

		reader := bufio.NewReader(workerConn)
		buffer, err := reader.ReadBytes('\n')
		if err != nil {
			log.Println(err)
			return
		}
		msg := string(buffer)
		tokens := strings.Fields(msg)

		if tokens[0] != "SORTED" {
			log.Println(errors.New("Received message different than SORTED"))
			return
		}

		sortedChunk, err := utils.UncompressChunk(tokens[1])

		ch <- sortedChunk
	}(ch, remoteIP)

	select {
	case sortedChunk := <-ch:
		utils.StoreChunk(sortedChunk)
	case <-time.After(5 * time.Second):
		ctx.Ch() <- chunkToSort
		return fmt.Errorf("TIMEOUT: Machine %s timeouted during sorting of chunk %d", remoteIP, chunkToSort.ID)
	}

	return nil
}
