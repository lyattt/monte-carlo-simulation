//simulates one trial
package main

import (
    "fmt"
    "math/rand"
)

func simulateOneTrial() float64 {
    value := 100.0
    mu := 0.001
    sigma := 0.02

    for i := 0; i < 252; i++ {
        r := mu + sigma*rand.NormFloat64()
        value = value * (1 + r)
    }

    return value
}

// worker funcition
func worker(id int, numTrials int, results chan<- float64) {
	for i := 0; i < numTrials; i++ {
		res := simulateOneTrial() 
		results <- res            
	}
}

// runs multiple trials conccurently and collects results
func runConcurrentTrials(totalTrials int, workers int) []float64 {
	results := make(chan float64)

	trialsPerWorker := totalTrials / workers
	remaining := totalTrials % workers

	// launch goroutines
	for w := 0; w < workers; w++ {
		num := trialsPerWorker
		if w == workers-1 {
			num += remaining // last worker takes remainder
		}
		go worker(w, num, results)
	}

	// collect results from all workers
	values := make([]float64, 0, totalTrials)
	for i := 0; i < totalTrials; i++ {
		val := <-results
		values = append(values, val)
	}

	return values
}

func main() {
    results := runConcurrentTrials(10000, 4)
    fmt.Println("Total:", len(results))
}