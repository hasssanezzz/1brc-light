package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

var chunksize int

type Result struct {
	Min, Max, Sum, Visits int
}

// var hashtable map[string]*Result

func work(index int, chunk string, wg *sync.WaitGroup) {

	fmt.Printf("thread[%d]:\n%s\n", index, chunk)
	wg.Done()
}

func read_chunk(index int, file *os.File) string {

	_, err := file.Seek(int64(index*chunksize), io.SeekStart)
	if err != nil {
		fmt.Printf("error in seeking 0: %v\n", err)
	}

	reader := bufio.NewReader(file)

	buff := make([]byte, chunksize)
	_, err = reader.Read(buff)
	if err != nil {
		log.Fatal("read_chunk init ", err)
	}

	chunk := string(buff)

	// cut the start
	if index != 0 {
		for i := 0; i < len(chunk); i++ {
			if chunk[i] == ';' {
				chunk = chunk[i+1:]
				break
			}
		}
	}

	// include the end
	extraChunk, err := reader.ReadSlice(byte(';'))
	if err != nil {
		log.Fatal("read_chunk: ", err)
	}

	chunk += string(extraChunk)

	return chunk
}

func initiate_workers(filename string, number_of_workers int) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln("initiate workers:", err)
	}
	defer file.Close()

	var wg sync.WaitGroup
	wg.Add(number_of_workers)

	for i := 0; i < number_of_workers; i++ {
		chunk := read_chunk(i, file)
		fmt.Print(chunk)
		go work(i, chunk, &wg)
	}

	wg.Wait()

	file.Close()
}

func main() {

	// ? optimal chunk size = sqrt(file_size)
	// ? make pages into account

	filename := os.Args[1]

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("open file: ", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		log.Fatal("read stats: ", err)
	}

	filesize := stat.Size()

	number_of_workers := 12
	chunksize = int(filesize) / number_of_workers
	initiate_workers(filename, number_of_workers)
}
