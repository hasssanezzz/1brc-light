package main

import (
	"strconv"
	"strings"
	"sync"
)

type Result struct {
	Min, Max, Sum, Visits int
}

type Entry struct {
	Gover    string
	Temp     int
	sourceId int
}

type Worker struct {
	Id       int
	RecvChan <-chan string
	RsltChan chan<- Entry
}

func NewWorker(id int, recvchan <-chan string, rsltChan chan<- Entry) *Worker {
	return &Worker{
		Id:       id,
		RecvChan: recvchan,
		RsltChan: rsltChan,
	}
}

func (w *Worker) Consume(wg *sync.WaitGroup) {
	for chunk := range w.RecvChan {
		// fmt.Printf("Worker %d recv a chunk of size %d\n", w.Id, len(chunk))

		entries := strings.Split(chunk, ";")
		for _, entry := range entries {
			gover, tempStr, ok := strings.Cut(entry, ",")
			if !ok {
				continue
			}
			temp, err := strconv.Atoi(tempStr)
			if err != nil {
				continue
			}

			w.RsltChan <- Entry{gover, temp, w.Id}
		}
	}
	wg.Done()
}
