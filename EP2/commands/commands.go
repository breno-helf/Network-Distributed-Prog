package commands

import (
	"errors"
	"fmt"
	"net"

	"../utils"
)

// ENTER command will allow someone to enter in the network
func ENTER(conn net.Conn, ctx utils.Context) error {
	if !ctx.IsMasterNode {
		return errors.New("Can't let someone enter if I am not master node")
	}

	remoteAddr := conn.RemoteAddr().String()
	ctx.Nodes = append(ctx.Nodes, remoteAddr)

	_, err := conn.Write([]byte(fmt.Sprintf("LEADER %s", ctx.Leader)))

	if err != nil {
		return err
	}

	return nil
}
