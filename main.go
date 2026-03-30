//simulates one trial
package main

import (
    "fmt"
    "math/rand"
    "math"
	"flag"
	"time"

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

// worker function
func worker(id int, numTrials int, results chan<- float64) {
	for i := 0; i < numTrials; i++ {
		res := simulateOneTrial() 
		results <- res         
	}
}

// runs multiple trials conccurently and collects results
func runConcurrentTrials(totalTrials int, workers int) []float64 {
	results := make(chan float64, totalTrials)

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
	trials := flag.Int("trials", 10000, "total number of simulations")
	workers := flag.Int("workers", 4, "number of goroutines")
	seed := flag.Int64("seed", 0, "random seed(optional)")
	flag.Parse()

	var actualSeed int64 
	if *seed == 0 {
		actualSeed = time.Now().UnixNano()
	} else {
		actualSeed = *seed
	}
	rand.Seed(actualSeed)
	start := time.Now()

    results := runConcurrentTrials(*trials, *workers)


	mean := 0.0
	for _, v := range results {
		mean += v
	}
	mean /= float64(*trials)

	var variance float64
	for _, v := range results {
		diff := v - mean
		variance += diff * diff
	}
	variance /= float64(*trials - 1)

	stdDev := math.Sqrt(variance)
	stdErr := stdDev / math.Sqrt(float64(*trials))
	ci95 := 1.96 * stdErr

    //Output
    fmt.Printf("Monte Carlo Simulation Results\n")
    fmt.Printf("Trials: %d\n", *trials)
    fmt.Printf("Workers: %d\n", *workers)

    if(*seed != 0) {
        fmt.Printf("Seed: %d\n", *seed)
    }
	elapsed := time.Since(start)
    fmt.Printf("Time to run: %.1f seconds\n", elapsed.Seconds())
    fmt.Printf("Mean outcome: %.2f\n", mean)
    fmt.Printf("Std dev: %.2f\n", stdDev)
    fmt.Printf("95% CI: [%.2f, %.2f]\n", mean-ci95, mean+ci95)

}
