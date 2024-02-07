package main

import (
	"fmt"
	"sync"
	"time"
)

func merge(channels ...<-chan interface{}) <-chan interface{} {
	var wg sync.WaitGroup
	merged := make(chan interface{})

	for _, chIt := range channels {
		wg.Add(1)

		go func(ch <-chan interface{}) {
			defer wg.Done()
			for val := range ch {
				fmt.Println(val)
				merged <- val
			}
		}(chIt)
	}

	go func() {
		wg.Wait()
		close(merged)
	}()

	return merged
}

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-merge(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	fmt.Printf("fone after %v", time.Since(start))

}
