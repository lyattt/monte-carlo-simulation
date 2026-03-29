//simulates one trial
package main

import (
    "fmt"
    "math/rand"
    "math"
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

func main() {
    nTrials := 10000
    results := make([]float64, nTrials)

    for i := 0; i < nTrials; i++ {
		results[i] = simulateOneTrial()
	}

	mean := 0.0
	for _, v := range results {
		mean += v
	}
	mean /= float64(nTrials)

	var variance float64
	for _, v := range results {
		diff := v - mean
		variance += diff * diff
	}
	variance /= float64(nTrials - 1)
	stdDev := math.Sqrt(variance)
	stdErr := stdDev / math.Sqrt(float64(nTrials))

	ci95 := 1.96 * stdErr

    fmt.Println("Trials: ", nTrials)
    fmt.Println("Mean Outcome: ", mean)
    fmt.Println("Std dev: ", stdDev)
    fmt.Println("95% CI: [", ci95, "]")

}
