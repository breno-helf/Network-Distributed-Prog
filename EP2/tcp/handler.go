package tcp

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
	"../utils"
)

func handleCommand(conn net.Conn, msg string, ctx *utils.Context) error {
	tokens := strings.Fields(msg)
	cmd := tokens[0]
	fmt.Println("Received command", cmd)
	switch cmd {
	case "ENTER":
		err := commands.ENTER(conn, ctx)
		if err != nil {
			return err
		}
	case "LEADER":
		if len(tokens) < 2 {
			return errors.New("LEADER command requires an argument")
		}
		err := commands.LEADER(conn, ctx, tokens[1])
		if err != nil {
			return err
		}
	case "SORT":
		if len(tokens) < 3 {
			return errors.New("SORT command requires 2 arguments")
		}
		err := commands.SORT(conn, ctx, tokens[1], tokens[2])
		if err != nil {
			return err
		}
	// case "WORK":
	// case "DIED":
	default:
		return fmt.Errorf("Can't handle message '%s'", msg)
	}

	return nil
}

// HandleConnection knows how to handle a connection with Handler port
func HandleConnection(conn net.Conn, ctx *utils.Context) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		buffer, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}
		msg := string(buffer)

		err = handleCommand(conn, msg, ctx)
		if err != nil {
			log.Println(err)
		}
	}
}
