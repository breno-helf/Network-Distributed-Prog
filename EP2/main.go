package main

/* Made by:
 * Breno Helfstein Moura - 9790972
 * Matheus Barcellos de Castro Cunha - 11208238
**/
import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"./eventlog"
	"./master"
	"./slave"
	"./utils"
)

func main() {
	// Means that I have the file to be sorted and have to initialize the connections
	masterNode := false
	var listFilename string

	rand.Seed(time.Now().UTC().UnixNano())
	addressConfText, err := ioutil.ReadFile("address.conf")
	if err != nil {
		log.Fatal(utils.MAINERROR, err)
	}

	initialMachine := strings.TrimSpace(string(addressConfText))
	myIP, err := utils.GetMyIP()
	if err != nil {
		log.Fatal(utils.MAINERROR, err)
	}

	if len(os.Args) > 1 && initialMachine == myIP {
		masterNode = true
		listFilename = os.Args[1]
	}

	for _, arg := range os.Args {
		if arg == "-d" {
			eventlog.ActivateLogMode()
		}
	}

	if masterNode {
		master.Master(listFilename, myIP)
	} else {
		slave.Slave(initialMachine, myIP)
	}
}
