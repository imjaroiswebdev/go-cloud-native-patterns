package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

func circuitBreakerDemo() {
	fmt.Println("Circuit Breaker Pattern Demo...")
	ctx := context.Background()

	errorProneFeatureWithCircuitBreaker := Breaker(errorProneFeature, 2)
	for i := 0; i < 50; i++ {
		res, err := errorProneFeatureWithCircuitBreaker(ctx)
		if err != nil {
			log.Printf("[ERROR] %v", err)
		}

		fmt.Println(res)
		time.Sleep(300 * time.Millisecond) // Wait introduced to mimic service latency and let the Circuit Breaker backoff mechanism do its job.
	}
}

// We beging by creating a Circuit type that specifies the signature of the
// function that's interacting with your databse or other upstream service.

type Circuit func(ctx context.Context) (string, error)

func errorProneFeature(ctx context.Context) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	randNum := rand.Intn(1000)
	shouldFail := randNum >= 400
	if shouldFail {
		return "", errors.New("Internal Error")
	}

	return fmt.Sprintf("%d", randNum), nil
}

// The Breaker function accepts any function that conforms to the Circuit type
// definition, and a unsigned integer representing the number of consecutive
// failures allowed before the circuit automatically opens.

func Breaker(circuit Circuit, failureThreshold uint) Circuit {
	var m sync.RWMutex
	consecutiveFailures := 0
	lastAttempt := time.Now()

	return func(ctx context.Context) (string, error) {
		m.RLock() // Establish a "read lock" for reading the shared resource `consecutiveFailures`

		d := consecutiveFailures - int(failureThreshold)

		if d >= 0 {
			shouldRetryAt := lastAttempt.Add(100 * time.Millisecond << d) // retry exponential backoff mechanism by binary shifting left the duration
			if !time.Now().After(shouldRetryAt) {
				m.RUnlock()
				return "", errors.New("service unreachable")
			}
		}

		m.RUnlock() // Release read lock on `consecutiveFailures`

		response, err := circuit(ctx) // Issue request proper

		m.Lock() // Lock around shared resources
		defer m.Unlock()

		lastAttempt = time.Now() // Record time of attempt

		if err != nil { // Circuit returned an error,
			consecutiveFailures++ // so we count the failure
			return response, err  // and return
		}

		consecutiveFailures = 0 // Reset failures counter

		return response, nil
	}
}
