package patterns

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

func RetryDemo() {
	fmt.Println("Retry Pattern Demo...")

	ctx := context.Background()

	emulateTransientErrorWithRetry := Retry(emulateTransientError, 5, 2*time.Second)
	res, err := emulateTransientErrorWithRetry(ctx)

	fmt.Println(res, err)
}

// We beging by creating an Effector type that specifies the signature of the
// function that's interacting with your databse or other upstream service.

type Effector func(context.Context) (string, error)

func Retry(effector Effector, retries int, delay time.Duration) Effector {
	return func(ctx context.Context) (string, error) {
		for r := 0; ; r++ {
			response, err := effector(ctx)
			if err == nil || r >= retries {
				return response, err
			}

			log.Printf("Attempt %d failed; retrying in %v", r+1, delay)

			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return "", ctx.Err()
			}
		}
	}
}

var count int

func emulateTransientError(ctx context.Context) (string, error) {
	count++

	if count <= 3 {
		return "intentional fail", errors.New("error")
	} else {
		return "success", nil
	}
}
