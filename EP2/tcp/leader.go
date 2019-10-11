package tcp

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

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
	lastRequest := make(map[string]time.Time)

	conn, err := net.Dial("tcp", ctx.MasterNode()+utils.HandlerPort)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("LEADER started connection", conn.LocalAddr(), conn.RemoteAddr())

	for ctx.IsLeader() {
		if idx == len(allNodes) {
			idx = 0
			allNodes = ctx.AllNodes()
			rand.Shuffle(len(allNodes), func(i, j int) { allNodes[i], allNodes[j] = allNodes[j], allNodes[i] })
		}

		err := askWorkForNode(conn, allNodes[idx], lastRequest)
		if err != nil {
			log.Println(err)
		}
		idx++
	}
}

func askWorkForNode(conn net.Conn, remoteIP string, lastRequest map[string]time.Time) error {
	_, ok := lastRequest[remoteIP]
	if ok {
		diff := time.Since(lastRequest[remoteIP])
		if diff < 2*time.Second {
			time.Sleep(time.Millisecond*5 - diff)
		}
	}
	lastRequest[remoteIP] = time.Now()
	_, err := fmt.Fprintf(conn, "WORK %s\n", remoteIP)
	return err
}
