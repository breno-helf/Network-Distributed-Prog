package slave

/* Made by:
 * Breno Helfstein Moura - 9790972
 * Matheus Barcellos de Castro Cunha - 11208238
**/
import (
	"fmt"
	"log"
	"net"

	"../eventlog"
	"../tcp"
	"../utils"
)

func enterNetwork(ctx *utils.Context) {
	err := utils.TryEnterNetwork(ctx)
	for err != nil {
		err = utils.TryEnterNetwork(ctx)
	}
	ctx.AddNode(ctx.MyIP())
	eventlog.EventNewNode(ctx.MyIP())
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
	go enterNetwork(ctx)

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
