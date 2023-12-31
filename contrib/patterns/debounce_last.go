package patterns

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

func DebounceLastDemo() {
	fmt.Println("Debounce Last Pattern Demo...")
	ctx := context.Background()

	statefulFeatureWithDebounceLast := DebounceLast(statefulFeature(), 100*time.Millisecond)
	for i := 0; i < 10; i++ {
		res, err := statefulFeatureWithDebounceLast(ctx)
		if err != nil {
			log.Printf("[ERROR] %v", err)
		}

		fmt.Println(res)

		var randNum float32 = 150 * rand.Float32()
		randLatency := time.Duration(randNum) * time.Millisecond
		time.Sleep(randLatency) // Wait introduced to mimic a random service latency
	}
}

// DebounceLast implementation involves the use of a `time.Ticker` to determine
// whether enough time has passed since the function was last called, and to
// call `circuit` when it has. Alternatively.
func DebounceLast(circuit Circuit, d time.Duration) Circuit {
	// This of `DebounceLast` takes pains to ensure thread safety by wrapping the
	// entire function in a mutex. While this will force overlapping calls at the
	// start of a cluster to have to wait until the result is cahed, it also
	// guarantees that `circuit` is called exactly once, at the very beginning of
	// a cluster.
	var m sync.Mutex
	var once sync.Once
	var threshold time.Time
	var ticker *time.Ticker

	var result string
	var err error

	return func(ctx context.Context) (string, error) {
		m.Lock()
		defer m.Unlock()

		threshold = time.Now().Add(d)

		// Almost the entire function is run inside this `Do` method of a
		// `sync.Once` value, which ensures that (as its name suggests) the
		// contained function is run exactly once.
		once.Do(func() {
			ticker = time.NewTicker(100 * time.Millisecond)

			go func() {
				defer func() {
					m.Lock()
					ticker.Stop()
					once = sync.Once{}
					m.Unlock()
				}()

				for {
					select {
					case <-ticker.C:
						m.Lock()
						if time.Now().After(threshold) {
							result, err = circuit(ctx)
							m.Unlock()
							return
						}
						m.Unlock()
					case <-ctx.Done():
						m.Lock()
						result, err = "", ctx.Err()
						m.Unlock()
						return
					}
				}
			}()

		})

		return result, err
	}
}
