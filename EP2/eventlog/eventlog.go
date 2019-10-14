package eventlog

import (
	"fmt"
	"log"
	"os"
)

const logFile = "eventLog.txt"

var logmode = false
var logfd *os.File
var eventLogger *log.Logger

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
		eventLogger = log.New(logfd, "Event Logger: ", log.LstdFlags)
		log.SetOutput(logfd)
		log.SetPrefix("Error Log: ")
		return err
	}

	return nil
}

// EventNewNode writes in the event logger that we have a new node
func EventNewNode(node string) {
	if logmode {
		eventLogger.Printf("New node entered the system [%s]\n", node)
	}
}

// EventDeadNode writes in the event logger that a node died
func EventDeadNode(node string) {
	if logmode {
		eventLogger.Printf("A node died [%s]\n", node)
	}
}

// EventElectingLeader writes in the event logger that we are electing a new leader
func EventElectingLeader() {
	if logmode {
		eventLogger.Print("Electing new leader\n")
	}
}

// EventLeaderElected writes in the event logger that we elected a new leader
func EventLeaderElected(node string) {
	if logmode {
		eventLogger.Printf("We have elected a new leader [%s]\n", node)
	}
}

// EventFinishSorting writes in the event logger that we finished sorting
func EventFinishSorting(masterNode string) {
	if logmode {
		eventLogger.Printf("We have finished sorting, we can find the array in [%s]\n", masterNode)
	}
}

// LogEvent log any event not specified
func LogEvent(msg string) {
	if logmode {
		eventLogger.Print(msg)
	}
}
