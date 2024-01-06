package patterns

import (
	"fmt"
	"sync"
	"time"
)

func FanInDemo() {
	fmt.Println("Fan-in Pattern Demo...")

	sources := make([]<-chan int, 0)

	for i := 0; i < 3; i++ {
		ch := make(chan int)
		sources = append(sources, ch)

		go func() {
			defer close(ch)

			for i := 1; i <= 5; i++ {
				ch <- i
				time.Sleep(100 * time.Millisecond)
			}
		}()
	}

	dest := Funnel(sources...)
	for d := range dest {
		fmt.Println(d)
	}
}

// Funnel is implemented as a function that receives zero to N input channels (Sources). For each input channel in Sources, the Funnel function starts a separate goroutine to read values from its assigned channel and forward them to a single output channel shared by all of the goroutines (Destination).
func Funnel(sources ...<-chan int) <-chan int {
	dest := make(chan int) // The shared output channel

	var wg sync.WaitGroup // Used to automatically close dest when all sources are closed

	wg.Add(len(sources)) // Set size of the WaitGroup

	for _, ch := range sources { // Start a goroutine for each source
		go func(c <-chan int) {
			defer wg.Done() // Notify WaiGroup when c closes

			for n := range c {
				dest <- n
			}
		}(ch)
	}

	go func() { // Start a goroutine to close dest after all sources close
		wg.Wait()
		close(dest)
	}()

	return dest
}
