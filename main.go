
// simulates one trial
package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"sync/atomic"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
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

// worker runs numTrials simulations and sends each result to the results channel.
// It also atomically increments a shared counter so the progress reporter can
// track how many trials have completed without blocking the workers.
func worker(id int, numTrials int, results chan<- float64, counter *int64) {
	for i := 0; i < numTrials; i++ {
		res := simulateOneTrial()
		results <- res
		atomic.AddInt64(counter, 1)
	}
}

// streamProgress prints a progress update every reportInterval completed trials.
// It runs in its own goroutine and stops when done is closed.
func streamProgress(counter *int64, totalTrials int, reportInterval int, done <-chan struct{}) {
	ticker := time.NewTicker(100 * time.Millisecond) // check counter 10x per second
	defer ticker.Stop()

	lastReported := 0

	for {
		select {
		case <-done:
			// Print a final 100% line when all trials finish
			fmt.Printf("\rProgress: %d / %d trials complete\n", totalTrials, totalTrials)
			return
		case <-ticker.C:
			completed := int(atomic.LoadInt64(counter))
			// Report only when we cross the next reportInterval boundary
			milestone := (completed / reportInterval) * reportInterval
			if milestone > lastReported && milestone > 0 {
				fmt.Printf("\rProgress: %d / %d trials complete", milestone, totalTrials)
				lastReported = milestone
			}
		}
	}
}

// runConcurrentTrials launches workers goroutines that together run totalTrials
// simulations. A background goroutine streams progress every reportInterval trials.
func runConcurrentTrials(totalTrials int, numWorkers int, reportInterval int) []float64 {
	results := make(chan float64, totalTrials)

	var counter int64 // atomically incremented by workers
	done := make(chan struct{})

	// Start the progress reporter before launching workers
	go streamProgress(&counter, totalTrials, reportInterval, done)

	trialsPerWorker := totalTrials / numWorkers
	remaining := totalTrials % numWorkers

	for w := 0; w < numWorkers; w++ {
		num := trialsPerWorker
		if w == numWorkers-1 {
			num += remaining // last worker takes the remainder
		}
		go worker(w, num, results, &counter)
	}

	// Collect every result from the channel
	values := make([]float64, 0, totalTrials)
	for i := 0; i < totalTrials; i++ {
		val := <-results
		values = append(values, val)
	}

	close(done) // signal the progress reporter to print its final line and exit
	return values
}

// saveHistogram builds a histogram of portfolio final values and writes it to
// a PNG file using gonum/plot.
func saveHistogram(values []float64, filename string) error {
	pts := make(plotter.Values, len(values))
	for i, v := range values {
		pts[i] = v
	}

	p := plot.New()
	p.Title.Text = "Monte Carlo: Final Portfolio Values"
	p.X.Label.Text = "Portfolio Value ($)"
	p.Y.Label.Text = "Count"

	// Number of bins follows a simple rule-of-thumb (square root of n, capped)
	numBins := int(math.Sqrt(float64(len(values))))
	if numBins > 200 {
		numBins = 200
	}
	if numBins < 10 {
		numBins = 10
	}

	h, err := plotter.NewHist(pts, numBins)
	if err != nil {
		return fmt.Errorf("creating histogram: %w", err)
	}

	p.Add(h)

	if err := p.Save(8*vg.Inch, 5*vg.Inch, filename); err != nil {
		return fmt.Errorf("saving histogram: %w", err)
	}
	return nil
}

func main() {
	trials  := flag.Int("trials", 10000, "total number of simulations")
	workers := flag.Int("workers", 4, "number of goroutines")
	seed    := flag.Int64("seed", 0, "random seed (optional)")
	flag.Parse()

	var actualSeed int64
	if *seed == 0 {
		actualSeed = time.Now().UnixNano()
	} else {
		actualSeed = *seed
	}
	rand.Seed(actualSeed) //nolint:staticcheck

	start := time.Now()

	// Stream a progress update every 10 000 trials
	reportInterval := 10000
	if *trials < reportInterval {
		reportInterval = *trials // avoid never printing anything for small runs
	}

	results := runConcurrentTrials(*trials, *workers, reportInterval)

	// --- Statistics ---
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
	stdErr  := stdDev / math.Sqrt(float64(*trials))
	ci95    := 1.96 * stdErr

	// --- Output ---
	fmt.Printf("\nMonte Carlo Simulation Results\n")
	fmt.Printf("Trials:       %d\n", *trials)
	fmt.Printf("Workers:      %d\n", *workers)
	if *seed != 0 {
		fmt.Printf("Seed:         %d\n", *seed)
	}
	elapsed := time.Since(start)
	fmt.Printf("Time to run:  %.1f seconds\n", elapsed.Seconds())
	fmt.Printf("Mean outcome: %.2f\n", mean)
	fmt.Printf("Std dev:      %.2f\n", stdDev)
	fmt.Printf("95%% CI:       [%.2f, %.2f]\n", mean-ci95, mean+ci95)

	// --- Histogram ---
	histFile := "portfolio_histogram.png"
	if err := saveHistogram(results, histFile); err != nil {
		fmt.Printf("Warning: could not save histogram: %v\n", err)
	} else {
		fmt.Printf("Histogram saved to %s\n", histFile)
	}
}
