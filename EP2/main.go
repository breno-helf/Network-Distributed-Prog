package main

/* Made by:
 * Breno Helfstein Moura - 9790972
 * Matheus Barcellos de Castro Cunha - 11208238
**/
import (
	"io/ioutil"
	"log"
	"os"

	"./eventlog"
	"./master"
	"./slave"
	"./utils"
)

func main() {
	// Means that I have the file to be sorted and have to initialize the connections
	masterNode := false
	var listFilename string

	initialMachine, err := ioutil.ReadFile("address.conf")
	if err != nil {
		log.Fatal(err)
	}

	myIP, err := utils.GetMyIP()
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 1 && string(initialMachine) == myIP {
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
		slave.Slave(string(initialMachine))
	}
}
