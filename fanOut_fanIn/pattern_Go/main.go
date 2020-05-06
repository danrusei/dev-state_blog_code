package main

import (
	"fmt"
	"math/rand"
	"time"
)

func generator(done <-chan struct{}) <-chan int {
	stream := make(chan int)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	i := 0

	go func() {
		defer close(stream)
		for {
			i = r.Intn(100)
			select {
			case <-done:
				return
			default:
				stream <- i
			}

		}
	}()
	return stream
}

func main() {

	done := make(chan struct{})
	i := 0

	for n := range generator(done) {
		if i > 10 {
			done <- struct{}{}
		}
		fmt.Println(n)
		i++
	}

}