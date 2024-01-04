# Contrib

Place meant to store Implementation of Cloud Native Patterns read from the [book](https://www.amazon.com/-/es/Matthew-Titmus/dp/1492076333/ref=tmm_pap_swatch_0?_encoding=UTF8&qid=&sr=) and their application, moreover, It will also be the place for examples of the main library.

## Patterns Implementation

### Stability patterns

* Circuit Breaker
* Debounce Function First
* Debounce Function Last
* Retry
* Throttle
* Timeout

### Patterns demo usage

> Providing no flag at all executes all the demos.

```sh
$ go run ./cmd/demo --help
Usage of /var/folders/.../demo:
  -circuit-breaker
        Execute Circuit Breaker Demo
  -debounce-first
        Execute Debounce First Demo
  -debounce-last
        Execute Debounce Last Demo
  -retry
        Execute Retry Demo
  -throttle
        Execute Throttle Demo
  -timeout
        Execute Time Demo
```

