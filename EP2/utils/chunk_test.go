package utils

/* Made by:
 * Breno Helfstein Moura - 9790972
 * Matheus Barcellos de Castro Cunha - 11208238
**/
import (
	"bufio"
	"os"
	"strconv"
	"testing"
)

func TestStoreChunk(t *testing.T) {
	type args struct {
		chunk Chunk
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		filename string
	}{
		{
			name:     "Testing storing a sorted chunk",
			args:     args{Chunk{[]int{1, 3, 4}, 2}},
			wantErr:  false,
			filename: "./tmp/chunk2",
		},
		{
			name:     "Testing storing an unsorted chunk",
			args:     args{Chunk{[]int{1, 3, 2}, 2}},
			wantErr:  true,
			filename: "./tmp/chunk2",
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storeErr := StoreChunk(tt.args.chunk)
			_, fileErr := os.Stat(tt.filename)
			if (storeErr != nil) != tt.wantErr {
				t.Errorf("storeChunk() error = %v, wantErr %v", storeErr, tt.wantErr)
			}

			if os.IsNotExist(fileErr) && tt.wantErr == false {
				t.Errorf("File was not created")
			} else {
				os.Remove(tt.filename)
			}

		})
	}

	os.RemoveAll("./tmp/")
}

func TestSortStoredChunks(t *testing.T) {
	chunks := []Chunk{
		Chunk{[]int{1, 5, 8, 11}, 0},
		Chunk{[]int{2, 3, 6, 22}, 1},
		Chunk{[]int{3, 4, 7}, 2},
	}

	for _, chunk := range chunks {
		StoreChunk(chunk)
	}

	if err := SortStoredChunks(len(chunks)); err != nil {
		t.Errorf("SortStoredChunks() error = %v", err)
	}

	sorted, err := os.Open(sortedFile)
	if err != nil {
		t.Errorf("Sorted file was not created %v", err)
	}
	defer sorted.Close()

	scanner := bufio.NewScanner(sorted)
	sortedSlice := []int{}

	for scanner.Scan() {
		number, err := strconv.Atoi(scanner.Text())
		if err != nil {
			t.Errorf("Couldnt convert string to int %v", err)
		}
		sortedSlice = append(sortedSlice, number)
	}

	CleanTemporaryFiles()
	os.Remove(sortedFile)
}
