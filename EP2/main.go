package main

/* Made by:
 * Breno Helfstein Moura - 9790972
 * Matheus Barcellos de Castro Cunha - 11208238
**/
import (
	"io/ioutil"
	"log"
	"os"

	"./utils"
)

func slave() {

}

func main() {
	// Means that I have the file to be sorted and have to initialize the connections
	initializer := false
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
		initializer = true
		listFilename = os.Args[1]
	}

	if initializer {
		master(listFilename, myIP)
	} else {
		slave()
	}
}
