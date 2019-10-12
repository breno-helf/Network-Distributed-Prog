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

func handleCommand(conn net.Conn, msg string, ctx *utils.Context, ch chan bool) error {
	tokens := strings.Fields(msg)
	cmd := tokens[0]
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
			log.Println(err)
		}
	case "SORT":
		if len(tokens) < 2 {
			return errors.New("SORT command requires an argument")
		}
		err := commands.SORT(conn, ctx, tokens[1])
		if err != nil {
			return err
		}
	case "WORK":
		if len(tokens) < 2 {
			return errors.New("WORK command requires an argument")
		}
		ch <- true
		go func(conn net.Conn, tokens []string, ctx *utils.Context) {
			err := commands.WORK(conn, ctx, tokens[1])
			if err != nil {
				log.Println(err)
			}
			<-ch
		}(conn, tokens, ctx)
	case "ENTERED":
		if len(tokens) < 2 {
			return errors.New("ENTERED command requires an argument")
		}
		err := commands.ENTERED(conn, ctx, tokens[1])
		if err != nil {
			return err
		}
	case "ELECTION":
		err := commands.ELECTION(conn, ctx)
		if err != nil {
			return err
		}
	case "END":
		err := commands.END(conn, ctx)
		if err != nil {
			return err
		}
	case "NODES":
		if len(tokens) < 2 {
			return errors.New("NODES command requires an argument")
		}
		err := commands.NODES(conn, ctx, tokens[1])
		if err != nil {
			return err
		}

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
	ch := make(chan bool, 5)
	defer close(ch)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		err = handleCommand(conn, msg, ctx, ch)
		if err != nil {
			log.Println(err)
		}
	}
}
