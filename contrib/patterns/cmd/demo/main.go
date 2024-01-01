package main

import (
	"flag"
	"fmt"

	patterns "github.com/imjaroiswebdev/go-cloud-native-patterns-contrib-patterns"
)

func main() {
	circuitBreakerFlag := flag.Bool("circuit-breaker", false, "Execute Circuit Breaker Demo")
	debounceFirstFlag := flag.Bool("debounce-first", false, "Execute Debounce First Demo")
	debounceLastFlag := flag.Bool("debounce-last", false, "Execute Debounce Last Demo")
	retryFlag := flag.Bool("retry", false, "Execute Retry Demo")

	flag.Parse()

	if *circuitBreakerFlag {
		patterns.CircuitBreakerDemo()
	}
	if *debounceFirstFlag {
		patterns.DebounceFirstDemo()
	}
	if *debounceLastFlag {
		patterns.DebounceLastDemo()
	}
	if *retryFlag {
		patterns.RetryDemo()
	}

	// If no flags are set, execute all demos
	if !(*circuitBreakerFlag || *debounceFirstFlag || *debounceLastFlag || *retryFlag) {
		fmt.Println("Executing all demos...")
		patterns.CircuitBreakerDemo()
		patterns.DebounceFirstDemo()
		patterns.DebounceLastDemo()
		patterns.RetryDemo()
	}

}
