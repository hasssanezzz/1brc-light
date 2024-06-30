package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

var chunksize int
var finalResults = make(map[string]*Result)
var chunkChan chan string
var entryChan chan Entry
var start = time.Now()

func read_chunks(file *os.File, streamChan chan<- string) {
	reader := bufio.NewReader(file)

loop:
	for {
		buff := make([]byte, chunksize)

		n, err := reader.Read(buff)
		switch err {
		case nil:
		case io.EOF:
			streamChan <- string(buff[:n])
			break loop
		default:
			fmt.Println("error in reading chunk:", err)
			return
		}

		if n != chunksize {
			buff = buff[:n]
		}

		for buff[len(buff)-1] != ';' {
			readbyte, err := reader.ReadByte()
			if err != nil {
				fmt.Printf("error in reading single byte: %v", err)
				break
			}

			buff = append(buff, readbyte)
		}

		streamChan <- string(buff)
	}
}

func recvAndWrite() {
	for entry := range entryChan {
		// fmt.Printf("[recv] (%s, %d) from %d\n", entry.Gover, entry.Temp, entry.sourceId)

		result, ok := finalResults[entry.Gover]
		if !ok {
			r := &Result{
				Min:    entry.Temp,
				Max:    entry.Temp,
				Sum:    entry.Temp,
				Visits: 1,
			}
			finalResults[entry.Gover] = r
		} else {
			result.Min = min(result.Min, entry.Temp)
			result.Max = max(result.Max, entry.Temp)
			result.Sum += entry.Temp
			result.Visits++
		}

	}

	println("results written:", time.Since(start).String())
}

func initiate_workers(filename string, number_of_workers int) {
	go func() {
		file, err := os.Open(filename)
		if err != nil {
			log.Fatalln("initiate workers:", err)
		}
		defer file.Close()

		read_chunks(file, chunkChan)
		close(chunkChan)
	}()

	var wg sync.WaitGroup
	wg.Add(number_of_workers)
	for i := 0; i < number_of_workers; i++ {
		w := NewWorker(i, chunkChan, entryChan)
		go w.Consume(&wg)
	}

	go recvAndWrite()

	// how TF does this work??????????
	wg.Wait()
	close(entryChan)
	println("entry channel closed:", time.Since(start).String())
}

func main() {

	// ? optimal chunk size = sqrt(file_size)

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

	number_of_workers := 100
	chunksize = int(filesize) / number_of_workers
	chunkChan = make(chan string, number_of_workers)
	entryChan = make(chan Entry, 1)
	initiate_workers(filename, number_of_workers)

	for gover, result := range finalResults {
		fmt.Printf("%s=%d/%d/%d\n", gover, result.Min, result.Max, result.Sum/result.Visits)
	}
}
