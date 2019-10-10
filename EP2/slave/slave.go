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

	"../tcp"
	"../utils"
)

func enterNetwork(ctx *utils.Context) error {
	conn, err := net.Dial("tcp", ctx.MasterNode()+utils.HandlerPort)
	if err != nil {
		return err
	}
	defer conn.Close()

	conn.Write([]byte("ENTER\n"))

	reader := bufio.NewReader(conn)
	msg, err := reader.ReadString('\n')
	tokens := strings.Fields(msg)

	if tokens[0] == "LEADER" {
		ctx.ChangeLeader(tokens[1])
		fmt.Println("Leader is", tokens[1])
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

// Slave defines the behaviour of a slave node
func Slave(masterNode string, myIP string) {
	ctx := utils.NewContext(
		[]string{masterNode},
		masterNode,
		masterNode,
		myIP,
		nil,
	)

	fmt.Println("Started slave")
	go messenger(ctx)

	listener, err := net.Listen("tcp", utils.HandlerPort)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Println("Slave server started! Waiting for connections...")

	connCh := utils.ClientConns(listener)

	for conn := range connCh {
		go tcp.HandleConnection(conn, ctx)
	}

}
