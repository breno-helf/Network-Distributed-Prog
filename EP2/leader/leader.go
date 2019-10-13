package leader

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"

	"../utils"
)

/* Made by:
 * Breno Helfstein Moura - 9790972
 * Matheus Barcellos de Castro Cunha - 11208238
**/

// Leader is for the node to act as Leader
func Leader(ctx *utils.Context) {
	idx := 0
	allNodes := ctx.AllNodes()
	rand.Shuffle(len(allNodes), func(i, j int) {
		allNodes[i], allNodes[j] = allNodes[j], allNodes[i]
	})

	conn, err := net.Dial("tcp", ctx.MasterNode()+utils.HandlerPort)
	if err != nil {
		log.Printf(utils.LEADERERROR, err)
		return
	}
	defer conn.Close()

	go checkWorkDone(conn, ctx)

	for ctx.IsLeader() {
		if idx == len(allNodes) {
			idx = 0
			allNodes = ctx.AllNodes()
			rand.Shuffle(len(allNodes), func(i, j int) { allNodes[i], allNodes[j] = allNodes[j], allNodes[i] })
		}

		err := askWorkForNode(conn, ctx, allNodes[idx])
		if err != nil {
			log.Printf(utils.LEADERERROR, err)
		}
		idx++
	}
}

func checkWorkDone(conn net.Conn, ctx *utils.Context) {
	reader := bufio.NewReader(conn)
	for ctx.IsLeader() {
		msg, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		tokens := strings.Fields(msg)
		if tokens[0] == "DONE" {
			err = DONE(conn, ctx, tokens[1])
			if err != nil {
				log.Printf(utils.LEADERERROR, err)
			}
		} else {
			log.Printf(utils.LEADERERROR, fmt.Errorf("Can't handle command %s on leader port", tokens[0]))
		}
	}
}

func askWorkForNode(conn net.Conn, ctx *utils.Context, remoteIP string) error {
	workLoad, ok := ctx.WorkLoad(remoteIP)
	if !ok {
		return fmt.Errorf("Node %s is not on the network anymore", remoteIP)
	}

	select {
	case workLoad <- true:
		_, err := fmt.Fprintf(conn, "WORK %s\n", remoteIP)
		return err
	default:
		// Node is with work load full
		return nil
	}
}

// DONE will notify master that a node finished sorting
func DONE(conn net.Conn, ctx *utils.Context, freeNode string) error {
	remoteIP := utils.GetRemoteIP(conn)
	if remoteIP != ctx.MasterNode() {
		return errors.New("Only master can inform of sorting done")
	}

	if !ctx.IsLeader() {
		return errors.New("Only leader can recognizes that a job is done")
	}

	workLoad, ok := ctx.WorkLoad(freeNode)
	if !ok {
		return fmt.Errorf("Node %s was removed from network", freeNode)
	}
	<-workLoad

	return nil
}
