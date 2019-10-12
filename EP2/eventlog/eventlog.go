package eventlog

import (
	"fmt"
	"log"
	"os"

	"../utils"
)

const logFile = "eventLog.txt"

var logmode = false
var logfd *os.File

// ActivateLogMode activates the log mode
func ActivateLogMode() {
	fmt.Println("LogMode activated")
	logmode = true
	err := CreateEventLogger()
	if err != nil {
		log.Printf("Failed to create logger: %v", err)
	}
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
		msg := fmt.Sprintf("New node entered the system [%s]\n", node)
		_, err := logfd.WriteString(msg)
		if err != nil {
			log.Println(utils.LOGERROR, err)
		}
	}
}

// EventDeadNode writes in the event logger that a node died
func EventDeadNode(node string) {
	if logmode {
		msg := fmt.Sprintf("A node died [%s]\n", node)
		_, err := logfd.WriteString(msg)
		if err != nil {
			log.Println(utils.LOGERROR, err)
		}
	}
}

// EventElectingLeader writes in the event logger that we are electing a new leader
func EventElectingLeader() {
	if logmode {
		_, err := logfd.WriteString(("Electing new leader"))
		if err != nil {
			log.Println(utils.LOGERROR, err)
		}
	}
}

// EventLeaderElected writes in the event logger that we elected a new leader
func EventLeaderElected(node string) {
	if logmode {
		msg := fmt.Sprintf("We have elected a new leader [%s]\n", node)
		_, err := logfd.WriteString(msg)
		if err != nil {
			log.Println(utils.LOGERROR, err)
		}
	}
}

// EventFinishSorting writes in the event logger that we finished sorting
func EventFinishSorting(masterNode string) {
	if logmode {
		msg := fmt.Sprintf("We have finished sorting, we can find the array in [%s]\n", masterNode)
		_, err := logfd.WriteString(msg)
		if err != nil {
			log.Println(utils.LOGERROR, err)
		}
	}
}
