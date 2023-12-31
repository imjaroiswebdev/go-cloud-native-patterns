package patterns

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

func DebounceFirstDemo() {
	fmt.Println("Debounce First Pattern Demo...")
	ctx := context.Background()

	statefulFeatureWithDebounceFirst := DebounceFirst(statefulFeature(), 100*time.Millisecond)
	for i := 0; i < 10; i++ {
		res, err := statefulFeatureWithDebounceFirst(ctx)
		if err != nil {
			log.Printf("[ERROR] %v", err)
		}

		fmt.Println(res)

		var randNum float32 = 150 * rand.Float32()
		randLatency := time.Duration(randNum) * time.Millisecond
		time.Sleep(randLatency) // Wait introduced to mimic a random service latency
	}
}

func statefulFeature() Circuit {
	count := 1

	return func(ctx context.Context) (string, error) {
		defer func() {
			count++
		}()

		return fmt.Sprintf("count = %d", count), nil
	}
}

// We start by defining a function type witht he signature of the function we
// want to limit. Also like Circuit Breaker, we call it `Circuit`. Due to it is
// identical to the one used in Circuit Breaker example, We're going to re-use
// it.

// The function-first implementation of Debounce (DebounceFirst) is very straigh
// forward compared to function-last because it only needs to track the last
// time it was called and return a cached result if it's called again less than
// `d` duration after.
func DebounceFirst(circuit Circuit, d time.Duration) Circuit {
	// This of `DebounceFirst` takes pains to ensure thread safety by wrapping the
	// entire function in a mutex. While this will force overlapping calls at the
	// start of a cluster to have to wait until the result is cahed, it also
	// guarantees that `circuit` is called exactly once, at the very beginning of
	// a cluster.
	var m sync.Mutex
	var threshold time.Time

	var result string
	var err error

	return func(ctx context.Context) (string, error) {
		m.Lock()

		defer func() {
			threshold = time.Now().Add(d)
			m.Unlock()
		}()

		if time.Now().Before(threshold) {
			return result, err
		}

		result, err = circuit(ctx)

		return result, err
	}
}
