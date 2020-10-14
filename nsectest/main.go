package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sync"
)

const concurrent = 50

func main() {

	// Command line parameters
	getConfig()

	// collect results
	var maccess = &sync.Mutex{}
	var measurements = make([]measurement, 0)

	// calc salt for hashing
	salt := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		log.Fatal(err)
	}

	// syncing of go routines
	var wg sync.WaitGroup
	var threads = make(chan string, concurrent)
	defer close(threads)

	for _, resolver := range resolverlist {
		wg.Add(1)
		go func(resolver Resolver) {
			defer wg.Done()
			threads <- "x"
			defer func() { _ = <-threads }()

			m, err := measure(resolver)
			if err == nil {
				maccess.Lock()
				measurements = append(measurements, m)
				maccess.Unlock()
			}
		}(resolver)
	}

	// wait for results to come in
	wg.Wait()
	log.Printf("Number of results: %d\n", len(measurements))

	// write results to files
	b, _ := json.Marshal(measurements)
	fmt.Println(bytes.NewBuffer(b).String())
}
