package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

var chunksize int

type Result struct {
	Min, Max, Sum, Visits int
}

var finalResults = make(map[string]*Result)

func work(chunk string, mainchan chan<- map[string]*Result, wg *sync.WaitGroup) {
	defer wg.Done()

	localMap := map[string]*Result{}

	entries := strings.Split(chunk, ";")
	for _, entry := range entries {
		splitResult := strings.Split(entry, ",")
		if len(splitResult) != 2 {
			continue
		}

		gover, tempStr := splitResult[0], splitResult[1]
		temp, err := strconv.Atoi(tempStr)
		if err != nil {
			log.Println("error parsing temp, entry:", entry)
			continue
		}

		result, ok := localMap[gover]
		if !ok {
			r := &Result{
				Min:    temp,
				Max:    temp,
				Sum:    temp,
				Visits: 1,
			}
			localMap[gover] = r
		} else {
			if temp < result.Min {
				result.Min = temp
			}
			if temp > result.Max {
				result.Max = temp
			}
			result.Sum += temp
			result.Visits++
		}
	}

	mainchan <- localMap
}

func read_chunks(file *os.File) []string {
	reader := bufio.NewReader(file)

	// TODO gotta be cleaned
	var chunks []string

	for {
		buff := make([]byte, chunksize)
		n, err := reader.Read(buff)
		if err == io.EOF {
			chunks = append(chunks, string(buff[:n]))
			break
		}
		if err != nil {
			fmt.Println("error in reading chunk:", err)
			return chunks
		}

		if n != chunksize {
			buff = buff[:n]
		}

		for buff[len(buff)-1] != ';' {
			readbyte, err := reader.ReadByte()
			if err != nil {
				fmt.Printf("error in reading byte: %v", err)
				break
			}

			buff = append(buff, readbyte)
		}

		chunks = append(chunks, string(buff))
	}

	return chunks
}

func initiate_workers(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln("initiate workers:", err)
	}
	defer file.Close()

	// TODO stream the chunks
	chunks := read_chunks(file)

	for _, c := range chunks {
		fmt.Print(c)
	}

	// var wg sync.WaitGroup
	// wg.Add(len(chunks))
	// resultChan := make(chan map[string]*Result)

	// for _, chunk := range chunks {
	// 	go work(chunk, resultChan, &wg)
	// }

	// go func() {
	// 	wg.Wait()
	// 	println("Channel closed")
	// 	close(resultChan)
	// }()

	// for mp := range resultChan {
	// 	for gover, result := range mp {
	// 		storedResult, ok := finalResults[gover]
	// 		if !ok {
	// 			finalResults[gover] = result
	// 		} else {
	// 			if storedResult.Min > result.Min {
	// 				storedResult.Min = result.Min
	// 			}

	// 			if storedResult.Max < result.Max {
	// 				storedResult.Max = result.Max
	// 			}

	// 			storedResult.Sum += result.Sum
	// 			storedResult.Visits += result.Visits
	// 		}
	// 	}
	// }
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

	number_of_workers := 10
	chunksize = int(filesize) / number_of_workers
	initiate_workers(filename)

	for gover, result := range finalResults {
		fmt.Printf("%s=%d/%d/%d\n", gover, result.Min, result.Max, result.Sum/result.Visits)
	}
}
