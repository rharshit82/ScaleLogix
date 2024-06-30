package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"time"
)

var payloads [][]byte
var payloadsDir string

func loadPayloads() error {
	for i := 1; i <= 2; i++ {
		filename := fmt.Sprintf("%s/schema%d.json", payloadsDir, i)
		payload, err := os.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("failed to load %s: %w", filename, err)
		}
		payloads = append(payloads, payload)
	}
	return nil
}
func randomPayload() []byte {
	fileContent := payloads[rand.Intn(len(payloads))] // Randomly select a file's content
	var objects []map[string]interface{}
	if err := json.Unmarshal(fileContent, &objects); err != nil {
		log.Fatal("Error unmarshaling payload:", err)
	}
	// Randomly select an object from the array
	randomObject := objects[rand.Intn(len(objects))]
	payload, err := json.Marshal(randomObject)
	if err != nil {
		log.Fatal("Error marshaling payload:", err)
	}
	fmt.Print("Payload: ", string(payload), "\n")
	return payload
}

var (
	reqs int
	max  int
)

func init() {
	flag.IntVar(&reqs, "reqs", 100, "Total requests")
	flag.IntVar(&max, "concurrent", 2, "Maximum concurrent requests")
}

type Response struct {
	*http.Response
	err error
}

// Dispatcher
func dispatcher(reqChan chan *http.Request) {
	defer close(reqChan)
	for i := 0; i < reqs; i++ {
		payload := randomPayload()
		req, err := http.NewRequest("POST", "http://34.100.240.235", bytes.NewBuffer(payload))
		if err != nil {
			log.Println(err)
		} else {
			req.Header.Set("Content-Type", "application/json")
		}
		reqChan <- req
	}
}

// Worker Pool
func workerPool(reqChan chan *http.Request, respChan chan Response) {
	t := &http.Transport{}
	for i := 0; i < max; i++ {
		go worker(t, reqChan, respChan)
	}
}

// Worker
func worker(t *http.Transport, reqChan chan *http.Request, respChan chan Response) {
	for req := range reqChan {
		resp, err := t.RoundTrip(req)
		r := Response{resp, err}
		respChan <- r
	}
}

// Consumer
func consumer(respChan chan Response) (int64, int64) {
	var (
		conns int64
		size  int64
	)
	for conns < int64(reqs) {
		select {
		case r, ok := <-respChan:
			if ok {
				if r.err != nil {
					log.Println("Request error:", r.err)
				} else {
					body, err := io.ReadAll(r.Body)
					if err != nil {
						log.Println("Error reading response body:", err)
					} else {
						fmt.Println("Response Body:", string(body)) // Log the response body
						size += r.ContentLength
					}
					if err := r.Body.Close(); err != nil {
						log.Println("Error closing response body:", err)
					}
				}
				conns++
			}
		}
	}
	return conns, size
}

func main() {
	flag.StringVar(&payloadsDir, "payloadsDir", "./payloads", "Directory where payload files are located")

	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	if err := loadPayloads(); err != nil {
		log.Fatalf("Failed to load payloads: %v", err)
	}
	reqChan := make(chan *http.Request)
	respChan := make(chan Response)
	start := time.Now()
	go dispatcher(reqChan)
	go workerPool(reqChan, respChan)
	conns, size := consumer(respChan)
	took := time.Since(start)
	ns := took.Nanoseconds()
	av := ns / conns

	average, err := time.ParseDuration(fmt.Sprintf("%d", av) + "ns")
	if err != nil {
		log.Println(err)
	}
	totalSeconds := took.Seconds()      
	rps := float64(reqs) / totalSeconds // Calculate requests per second

	fmt.Printf("Connections:\t%d\nConcurrent:\t%d\nTotal size:\t%d bytes\nTotal time:\t%s\nAverage time:\t%s\nRequests per second:\t%.2f\n", conns, max, size, took, average, rps)
}
