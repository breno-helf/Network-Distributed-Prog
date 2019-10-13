package slave

/* Made by:
 * Breno Helfstein Moura - 9790972
 * Matheus Barcellos de Castro Cunha - 11208238
**/
import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"

	"../eventlog"
	"../tcp"
	"../utils"
)

func tryEnterNetwork(ctx *utils.Context) error {
	conn, err := net.Dial("tcp", ctx.MasterNode()+utils.HandlerPort)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = fmt.Fprintf(conn, "ENTER\n")
	if err != nil {
		return err
	}

	reader := bufio.NewReader(conn)
	msg, err := reader.ReadString('\n')
	tokens := strings.Fields(msg)

	if tokens[0] == "LEADER" {
		ctx.SetLeader(tokens[1])
	} else {
		return errors.New("Expecting LEADER message")
	}

	return nil
}

// Slave defines the behaviour of a slave node
func Slave(masterNode string, myIP string) {
	ctx := utils.NewContext(
		map[string]bool{masterNode: true},
		masterNode,
		masterNode,
		myIP,
		nil,
	)

	fmt.Println("Started slave")

	// Keep trying to connect
	go func(ctx *utils.Context) {
		err := tryEnterNetwork(ctx)
		for err != nil {
			err = tryEnterNetwork(ctx)
		}
		ctx.AddNode(ctx.MyIP())
		eventlog.EventNewNode(ctx.MyIP())
	}(ctx)

	listener, err := net.Listen("tcp", utils.HandlerPort)
	if err != nil {
		log.Fatal(utils.SLAVEERROR, err)
	}
	defer listener.Close()

	fmt.Println("Slave server started! Waiting for connections...")

	connCh := utils.ClientConns(listener)

	for conn := range connCh {
		go tcp.HandleConnection(conn, ctx)
	}

}
