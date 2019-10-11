package master

/* Made by:
 * Breno Helfstein Moura - 9790972
 * Matheus Barcellos de Castro Cunha - 11208238
**/
import (
	"bufio"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"strconv"
	"time"

	"../tcp"
	"../utils"
)

func generateChunks(filename string, ch chan<- utils.Chunk, chunkSize int) {
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
			log.Printf("String can't be properly converted to integer in list file: %v\n", err)
		}
		currentSlice = append(currentSlice, num)
		if len(currentSlice) == chunkSize {
			currentChunk := utils.Chunk{
				Numbers: currentSlice,
				ID:      currentID,
			}
			ch <- currentChunk
			currentID++
			currentSlice = nil
		}
	}

	if len(currentSlice) > 0 {
		currentChunk := utils.Chunk{
			Numbers: currentSlice,
			ID:      currentID,
		}
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

func election(ctx *utils.Context) {

}

func keepElecting(ctx *utils.Context) {
	for {
		time.Sleep(time.Minute)
		election(ctx)
	}
}

func waitForFinalSort(ctx *utils.Context, maxChunk int) {
	for !ctx.FinalSort() {
		time.Sleep(time.Second)
	}
	err := utils.SortStoredChunks(maxChunk)
	if err != nil {
		log.Fatal(err)
	}
	utils.CleanTemporaryFiles()
	os.Exit(0)
}

// Master executes the behaviour of the master node
func Master(listFilename string, myIP string) {
	// This channel will carry the chunks that will be sent to the other computers
	chunksChannel := make(chan utils.Chunk, 10)
	lineNumber := countLines(listFilename)
	chunkSize := defineChunkSize(lineNumber)
	numChunks := int(math.Ceil(float64(lineNumber) / float64(chunkSize)))
	ctx := utils.NewContext(map[string]bool{myIP: true}, myIP, myIP, myIP, chunksChannel)

	go generateChunks(listFilename, chunksChannel, chunkSize)

	listener, err := net.Listen("tcp", utils.HandlerPort)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	go tcp.Leader(ctx)
	go keepElecting(ctx)
	go waitForFinalSort(ctx, numChunks)

	fmt.Println("Master server started! Waiting for connections...")

	connCh := utils.ClientConns(listener)

	for conn := range connCh {
		fmt.Println("MASTER Started connection", conn.LocalAddr(), conn.RemoteAddr())
		go tcp.HandleConnection(conn, ctx)
	}
}
