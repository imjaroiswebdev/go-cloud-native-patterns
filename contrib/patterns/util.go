package patterns

import (
	"fmt"
	"time"
)

// spinner will print an spinner
func spinner(delay time.Duration) {
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}
