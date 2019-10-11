package utils

import (
	"bufio"
	"container/heap"
	"errors"
	"os"
	"sort"
	"strconv"
)

// Chunk is the chunk to be ordered by some node
type Chunk struct {
	Numbers []int
	ID      int
}

type fileChunk struct {
	scanner *bufio.Scanner
	head    int
	index   int
}

type fileChunkPQ []*fileChunk

func (pq fileChunkPQ) Len() int { return len(pq) }

func (pq fileChunkPQ) Less(i, j int) bool {
	return pq[i].head < pq[j].head
}

func (pq *fileChunkPQ) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

func (pq *fileChunkPQ) Push(x interface{}) {
	n := len(*pq)
	item := x.(*fileChunk)
	item.index = n
	*pq = append(*pq, item)
}

func (pq fileChunkPQ) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

const tmpFolder = "./tmp/"
const sortedFile = "./sorted.txt"

// StoreChunk stores a chunk. Returns an error if the chunk is not sorted.
func StoreChunk(chunk Chunk) error {
	if !sort.IntsAreSorted(chunk.Numbers) {
		return errors.New("Can only store sorted chunks")
	}

	_ = os.Mkdir(tmpFolder, os.ModePerm)
	chunkFile, err := os.Create(tmpFolder + "chunk" + strconv.Itoa(chunk.ID))
	if err != nil {
		return err
	}
	defer chunkFile.Close()

	for _, v := range chunk.Numbers {
		_, err = chunkFile.WriteString(strconv.Itoa(v) + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

// SortStoredChunks sort the stored chunks in a final file
func SortStoredChunks(maxChunk int) error {
	pq := make(fileChunkPQ, maxChunk)

	for i := 0; i < maxChunk; i++ {
		fd, err := os.Open(tmpFolder + "chunk" + strconv.Itoa(i))
		defer fd.Close()
		if err != nil {
			return err
		}

		scanner := bufio.NewScanner(fd)
		if !scanner.Scan() {
			continue
		}

		head, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return err
		}

		pq[i] = &fileChunk{
			scanner: scanner,
			head:    head,
			index:   i,
		}
	}
	heap.Init(&pq)
	sorted, err := os.Create(sortedFile)
	defer sorted.Close()

	for pq.Len() > 0 {
		fc := heap.Pop(&pq).(*fileChunk)
		sorted.WriteString(strconv.Itoa(fc.head) + "\n")
		if fc.scanner.Scan() {
			fc.head, err = strconv.Atoi(fc.scanner.Text())
			if err != nil {
				return err
			}
			heap.Push(&pq, fc)
		}
	}

	return nil
}

// CleanTemporaryFiles Clean all temporary files
func CleanTemporaryFiles() error {
	return os.RemoveAll(tmpFolder)
}
