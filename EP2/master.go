package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"./tcp"
)

type chunk struct {
	numbers []int
	ID      int
}

func generateChunks(filename string, ch chan<- chunk, chunkSize int) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fscanner := bufio.NewScanner(f)
	currentSlice := make([]int, 0)
	currentID := 0
	for fscanner.Scan() {
		num, err := strconv.Atoi(fscanner.Text())
		if err != nil {
			fmt.Printf("String can't be properly converted to integer in list file: %v", err)
		}
		currentSlice = append(currentSlice, num)
		if len(currentSlice) == chunkSize {
			currentChunk := chunk{currentSlice, currentID}
			ch <- currentChunk
			currentID++
			currentSlice = nil
		}
	}

	if len(currentSlice) > 0 {
		currentChunk := chunk{currentSlice, currentID}
		ch <- currentChunk
		currentID++
	}

	close(ch)
}

func countLines(filename string) int {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	count := 0
	fscanner := bufio.NewScanner(f)
	for fscanner.Scan() {
		fscanner.Text()
		count++
	}

	return count
}

const maxChunkSize int = 1000000
const minChunkSize int = 10

func defineChunkSize(lineNumber int) int {
	if lineNumber/100 > maxChunkSize {
		return maxChunkSize
	}

	if lineNumber/100 < minChunkSize {
		return minChunkSize
	}

	return lineNumber / 100
}

func master(listFilename string) {
	// This channel will carry the chunks that will be sent to the other computers
	chunksChannel := make(chan chunk, 10)
	lineNumber := countLines(listFilename)
	chunkSize := defineChunkSize(lineNumber)

	go generateChunks(listFilename, chunksChannel, chunkSize)

	server, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	fmt.Println("Server started! Waiting for connections...")
	connection, err := server.Accept()

	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	fmt.Println("Client connected")

	tf, _ := os.Open("lista.txt")
	tcp.SendChunk(connection, tf, 5)
}
