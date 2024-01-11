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
	throttleFlag := flag.Bool("throttle", false, "Execute Throttle Demo")
	timeoutFlag := flag.Bool("timeout", false, "Execute Time Demo")
	faninFlag := flag.Bool("fanin", false, "Execute Fan-in Demo")
	fanoutFlag := flag.Bool("fanout", false, "Execute Fan-out Demo")
	futureFlag := flag.Bool("future", false, "Execute Future Demo")

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
	if *throttleFlag {
		patterns.ThrottleDemo()
	}
	if *timeoutFlag {
		patterns.TimeoutDemo()
	}
	if *faninFlag {
		patterns.FaninnDemo()
	}
	if *fanoutFlag {
		patterns.FanoutDemo()
	}
	if *futureFlag {
		patterns.FutureDemo()
	}

	// If no flags are set, execute all demos
	if !(*circuitBreakerFlag ||
		*debounceFirstFlag ||
		*debounceLastFlag ||
		*retryFlag ||
		*throttleFlag ||
		*timeoutFlag ||
		*faninFlag ||
		*fanoutFlag ||
		*futureFlag) {
		fmt.Println("Executing all demos...")
		patterns.CircuitBreakerDemo()
		patterns.DebounceFirstDemo()
		patterns.DebounceLastDemo()
		patterns.RetryDemo()
		patterns.ThrottleDemo()
		patterns.TimeoutDemo()
		patterns.FaninnDemo()
		patterns.FanoutDemo()
		patterns.FutureDemo()
	}

}
