package patterns

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

func ThrottleDemo() {
	ctx := context.Background()

	dummyEffectorWithThrottle := Throttle(dummyEffector, 3, 1, 10*time.Millisecond)
	for i := 0; i < 20; i++ {
		res, err := dummyEffectorWithThrottle(ctx)
		if err != nil {
			log.Printf("[ERROR] %v", err)
		}

		fmt.Println(res)
		var randNum float32 = 10 * rand.Float32()
		randLatency := time.Duration(randNum) * time.Millisecond
		time.Sleep(randLatency) // Wait introduced to mimic a random service latency
	}
}

// Throttle implementation uses the most common algorithm for implementing
// rate-limiting behaviour, The Token Bucket (https://oreil.ly/5A5aP), which
// uses the analogy of a bucket that can hold some maximum number of tokens.
// When a function is called, a token is taken fromt he bucket, which then
// refills at some fixed rate.
func Throttle(e Effector, max uint, refill uint, d time.Duration) Effector {
	var tokens = max
	var once sync.Once

	return func(ctx context.Context) (string, error) {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}

		// Token refill logic which happens at a `d` rate.
		once.Do(func() {
			ticker := time.NewTicker(d)

			go func() {
				defer ticker.Stop()

				for {
					select {
					case <-ctx.Done():
						return

					case <-ticker.C:
						t := tokens + refill
						if t > max {
							t = max
						}
						tokens = t
					}
				}
			}()
		})

		if tokens <= 0 {
			return "", errors.New("too many calls")
		}

		tokens--

		return e(ctx)
	}
}

func dummyEffector(ctx context.Context) (string, error) {
	return "success", nil
}
