package commands

import (
	"errors"
	"fmt"
	"net"

	"../eventlog"
	"../utils"
)

// ENTER command will allow someone to enter in the network
func ENTER(conn net.Conn, ctx *utils.Context) error {
	if !ctx.IsMasterNode() {
		return errors.New("Can't let someone enter if I am not the master node")
	}

	remoteAddr := conn.RemoteAddr().String()
	ctx.AddNode(remoteAddr)

	_, err := conn.Write([]byte(fmt.Sprintf("LEADER %s\n", ctx.Leader)))

	if err != nil {
		return err
	}

	eventlog.EventNewNode(remoteAddr)

	return nil
}
