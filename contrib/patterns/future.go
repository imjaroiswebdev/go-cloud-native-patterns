package patterns

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func FutureDemo() {
	fmt.Println("Future Pattern Demo...")

	ctx := context.Background()
	future := timeConsumingFunction(ctx)

	res, err := future.Result()
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println(res)
}

// We declare Future which is the interface InnerFuture will satisfy

type Future interface {
	Result() (string, error)
}

// InnerFuture has one or more methods that satisfy the Future interface, which
// retrieve the values returned by the core function from the channels, cache
// them, and return them. If the values aren't available ont eh channel, the
// request blocks. If they have already been retrived, the cached values are
// returned.
type InnerFuture struct {
	once sync.Once
	wg   sync.WaitGroup

	res   string
	err   error
	resCh <-chan string
	errCh <-chan error
}

func (f *InnerFuture) Result() (string, error) {
	f.once.Do(func() {
		f.wg.Add(1)
		defer f.wg.Done()
		f.res = <-f.resCh
		f.err = <-f.errCh
	})

	f.wg.Wait()

	return f.res, f.err
}

// timeConsumingFunction is a SlowFunction.
// SlowFunction is a wrapper around the core functionality that you want to run
// concurrently. It has the job of creating thr results channels, running the
// core function in a goroutine, and creating and returning the Future
// implementation. (InnerFuture, in this case)
func timeConsumingFunction(ctx context.Context) Future {
	resCh := make(chan string)
	errCh := make(chan error)

	go func() {
		select {
		case <-time.After(2 * time.Second):
			resCh <- "I slept for 2 seconds"
			errCh <- nil
		case <-ctx.Done():
			resCh <- ""
			errCh <- ctx.Err()
		}
	}()

	return &InnerFuture{resCh: resCh, errCh: errCh}
}
