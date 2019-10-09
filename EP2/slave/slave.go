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

	"../commands"
	"../tcp"
	"../utils"
)

func enterNetwork(ctx *utils.Context) error {
	conn, err := net.Dial("tcp", ctx.MasterNode()+utils.HandlerPort)
	if err != nil {
		return err
	}

	conn.Write([]byte("ENTER\n"))

	reader := bufio.NewReader(conn)
	msg, err := reader.ReadString('\n')
	tokens := strings.Fields(msg)

	if tokens[0] == "LEADER" {
		ctx.ChangeLeader(tokens[1])
	} else {
		return errors.New("Expecting LEADER message")
	}

	return nil
}

func messenger(ctx *utils.Context) {
	err := enterNetwork(ctx)
	// Keep trying to connect
	for err != nil {
		err = enterNetwork(ctx)
	}

}

func handleCommand(conn net.Conn, msg string, ctx *utils.Context) error {
	tokens := strings.Fields(msg)
	cmd := tokens[0]
	switch cmd {
	case "ENTER":
		err := commands.ENTER(conn, ctx)
		if err != nil {
			return err
		}
	// case "LEADER":
	// case "WORK":
	// case "SORT":
	// case "DIED":
	default:
		return fmt.Errorf("Can't handle message '%s'", msg)
	}

	return nil
}

func handleConnection(conn net.Conn, ctx *utils.Context) {
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		err = handleCommand(conn, msg, ctx)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// Slave defines the behaviour of a slave node
func Slave(masterNode string) {
	ctx := utils.NewContext(
		false,
		false,
		[]string{masterNode},
		masterNode,
		masterNode,
		nil,
	)

	listener, err := net.Listen("tcp", utils.HandlerPort)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Println("Slave server started! Waiting for connections...")

	connCh := tcp.ClientConns(listener)

	for conn := range connCh {
		go handleConnection(conn, ctx)
	}

}
