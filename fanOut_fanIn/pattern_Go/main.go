package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

//generator returns a random number in the range of 1 to 100
func generator(done <-chan struct{}, iter int) <-chan int {
	streamID := make(chan int)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	genID := 0

	go func() {
		defer close(streamID)
		for i := 0; i < iter; i++ {
			genID = r.Intn(100)
			select {
			case <-done:
				return
			case streamID <- genID:
			}
		}
	}()
	return streamID
}

//fanOutFunc returns in a channel the result of the called function.
func fanOutFunc(done chan struct{}, in <-chan int) <-chan string {
	resultValue := make(chan string)
	go func() {
		defer close(resultValue)
		for n := range in {
			select {
			case <-done:
				log.Println("funOutFunc has been canceled")
				return
			case resultValue <- parse(done, n) + " _ Processed":
			}
		}
	}()
	return resultValue
}

//fanIn consolidate all the worker channels within one channel
func fanIn(done <-chan struct{}, cs ...<-chan string) <-chan string {
	var wg sync.WaitGroup
	resultValue := make(chan string)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	multiplex := func(c <-chan string) {
		defer wg.Done()
		for text := range c {
			select {
			case <-done:
				log.Println("funIn has been canceled")
				return
			case resultValue <- text:
			}
		}
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go multiplex(c)
	}

	//a goroutine to close out the "value" channel once all goroutines are done
	go func() {
		wg.Wait()
		close(resultValue)
	}()

	return resultValue
}

func main() {
	nWorkers := 4
	nJobs := 8

	done := make(chan struct{})
	defer close(done)

	fanOut := make([]<-chan string, nWorkers)

	// fanOut defines the number of concurrent "fanoutFunc" functions (goroutines)
	for i := 0; i < nWorkers; i++ {
		fanOut[i] = fanOutFunc(done, generator(done, nJobs/nWorkers))
	}

	//this pipeline yields the result of each channel of the fanOut slice
	for result := range fanIn(done, fanOut...) {
		fmt.Println(result)
	}

}
