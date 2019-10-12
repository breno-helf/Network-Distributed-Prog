package master

/* Made by:
 * Breno Helfstein Moura - 9790972
 * Matheus Barcellos de Castro Cunha - 11208238
**/
import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"../leader"
	"../tcp"
	"../utils"
)

func generateChunks(ctx *utils.Context, filename string, ch chan<- utils.Chunk, chunkSize int) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(utils.MASTERERROR, err)
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
			ctx.Wg().Add(1)
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
		ctx.Wg().Add(1)
	}

	close(ch)
}

func countLines(filename string) int {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(utils.MASTERERROR, err)
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
	var wg sync.WaitGroup
	nodes := ctx.AllNodes()
	ch := make(chan bool, 5)
	votes := make(map[string]int)
	var mu sync.Mutex

	for _, node := range nodes {
		ch <- true
		wg.Add(1)
		go func(ch <-chan bool, remoteIP string) {
			conn, err := net.Dial("tcp", remoteIP+utils.HandlerPort)
			if err != nil {
				log.Println(utils.ELECTIONERROR, err)
			}
			defer conn.Close()

			fmt.Fprint(conn, "ELECTION")

			reader := bufio.NewReader(conn)
			msg, err := reader.ReadString('\n')
			if err != nil {
				log.Println(utils.ELECTIONERROR, err)
			}

			tokens := strings.Fields(msg)
			if len(tokens) < 2 {
				log.Println(utils.ELECTIONERROR, errors.New("Can't cast vote"))
			}

			mu.Lock()
			votes[tokens[1]]++
			mu.Unlock()

			<-ch
			wg.Done()
		}(ch, node)
	}
	wg.Wait()

	president, maxVotes := ctx.MasterNode(), -1
	for k, v := range votes {
		if v > maxVotes {
			president = k
			maxVotes = v
		}
	}

	utils.Broadcast(ctx, fmt.Sprintf("LEADER %s\n", president))
}

func keepElecting(ctx *utils.Context) {
	for {
		time.Sleep(time.Minute)
		election(ctx)
	}
}

func waitForFinalSort(ctx *utils.Context, maxChunk int) {
	ctx.Wg().Wait()
	err := utils.SortStoredChunks(maxChunk)
	if err != nil {
		log.Fatal(utils.MASTERERROR, err)
	}
	utils.CleanTemporaryFiles()
	utils.Broadcast(ctx, "END\n")
}

// Master executes the behaviour of the master node
func Master(listFilename string, myIP string) {
	// This channel will carry the chunks that will be sent to the other computers
	chunksChannel := make(chan utils.Chunk, 10)
	lineNumber := countLines(listFilename)
	chunkSize := defineChunkSize(lineNumber)
	numChunks := int(math.Ceil(float64(lineNumber) / float64(chunkSize)))
	ctx := utils.NewContext(map[string]bool{myIP: true}, myIP, myIP, myIP, chunksChannel)

	go generateChunks(ctx, listFilename, chunksChannel, chunkSize)

	listener, err := net.Listen("tcp", utils.HandlerPort)
	if err != nil {
		log.Fatal(utils.MASTERERROR, err)
	}
	defer listener.Close()

	go leader.Leader(ctx)
	go keepElecting(ctx)
	go waitForFinalSort(ctx, numChunks)

	fmt.Println("Master server started! Waiting for connections...")

	connCh := utils.ClientConns(listener)

	for conn := range connCh {
		go tcp.HandleConnection(conn, ctx)
	}
}
