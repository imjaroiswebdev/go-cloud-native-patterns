package patterns

import (
	"context"
	"fmt"
	"time"
)

func TimeoutDemo() {
	fmt.Println("Timeout Pattern Demo...")

	ctx := context.Background()
	ctxt, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	fmt.Println("Waiting response from long running function...")
	go spinner(100 * time.Millisecond)
	dummySlowFunctionWithTimeout := Timeout(dummySlowFunction)
	res, err := dummySlowFunctionWithTimeout(ctxt, "slowly returned")

	fmt.Printf("\rres %q; err %v\n", res, err)
}

// We beging by creating an SlowFunction type that specifies the signature of the
// long running function.

type SlowFunction func(string) (string, error)

type WithContext func(context.Context, string) (string, error)

// This pattern is particularly useful whenever the function it's trying to be
// invoke is from third party dependencies and its signature is locked to be
// updated in order to accept a `context.Context`.
func Timeout(f SlowFunction) WithContext {
	return func(ctx context.Context, arg string) (string, error) {
		chres := make(chan string)
		cherr := make(chan error)

		go func() {
			res, err := f(arg)
			chres <- res
			cherr <- err
		}()

		select {
		case res := <-chres:
			return res, <-cherr
		case <-ctx.Done():
			return "", ctx.Err()

			// Although it's usually preferred to implement service timeouts using
			// context.Context, channel timeouts can also be implemented using the
			// channel provided by the `time.After` function. Like follows...
			// case <-time.After(10 * time.Second):
			// 	return "", errors.New("timed out")
		}
	}
}

func dummySlowFunction(s string) (string, error) {
	time.Sleep(10 * time.Second)
	return s, nil
}
