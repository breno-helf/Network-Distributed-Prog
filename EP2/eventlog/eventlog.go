package eventlog

import (
	"fmt"
	"log"
	"os"
)

const logFile = "eventLog.txt"

var logmode = false
var logfd *os.File

// ActivateLogMode activates the log mode
func ActivateLogMode() {
	logmode = true
}

// CreateEventLogger creates an event logger
func CreateEventLogger() error {
	if logmode {
		var err error
		logfd, err = os.Create(logFile)
		return err
	}

	return nil
}

// EventNewNode writes in the event logger that we have a new node
func EventNewNode(node string) {
	if logmode {
		msg := fmt.Sprintf("New node entered the system [%s]", node)
		_, err := logfd.Write([]byte(msg))
		if err != nil {
			log.Println(err)
		}
	}
}

// EventDeadNode writes in the event logger that a node died
func EventDeadNode(node string) {
	if logmode {
		msg := fmt.Sprintf("A node died [%s]", node)
		_, err := logfd.Write([]byte(msg))
		if err != nil {
			log.Println(err)
		}
	}
}

// EventElectingLeader writes in the event logger that we are electing a new leader
func EventElectingLeader() {
	if logmode {
		_, err := logfd.Write([]byte("Electing new leader"))
		if err != nil {
			log.Println(err)
		}
	}
}

// EventLeaderElected writes in the event logger that we elected a new leader
func EventLeaderElected(node string) {
	if logmode {
		msg := fmt.Sprintf("We have elected a new leader [%s]", node)
		_, err := logfd.Write([]byte(msg))
		if err != nil {
			log.Println(err)
		}
	}
}

// EventFinishSorting writes in the event logger that we finished sorting
func EventFinishSorting(masterNode string) {
	if logmode {
		msg := fmt.Sprintf("We have finished sorting, we can find the array in [%s]", masterNode)
		_, err := logfd.Write([]byte(msg))
		if err != nil {
			log.Println(err)
		}
	}
}
