package patterns

import (
	"fmt"
	"sync"
)

func FanoutDemo() {
	fmt.Println("Retry Pattern Demo...")

	source := make(chan int)
	dests := Split(source, 5)

	go func() {
		for i := 1; i <= 10; i++ {
			source <- i
		}

		close(source)
	}()

	var wg sync.WaitGroup
	wg.Add(len(dests))

	for i, ch := range dests {
		go func(i int, d <-chan int) {
			defer wg.Done()

			for val := range d {
				fmt.Printf("Destination channel #%d got %d\n", i, val)
			}
		}(i, ch)
	}

	wg.Wait()
}

// Fan-out is implemented as a Split function, which accepts a single Source
// channel and an integer representing the desired number of Destination
// channels. The Split function creates the Destination channels and executes
// some background process that retrieves values from Source channel and
// forwards them to one of the Destinations.
func Split(source <-chan int, n int) []<-chan int {
	dests := make([]<-chan int, 0)

	// It will create separate goroutines for each Destiniation that compete to
	// read the next value from Source and forward it to their respective
	// Destination.
	for i := 0; i < n; i++ { // Create n destination channels
		ch := make(chan int)
		dests = append(dests, ch)

		go func() { // Each channel gets a dedicated
			defer close(ch) // goroutine that competes for reads

			for val := range source {
				ch <- val
			}
		}()
	}

	return dests
}
