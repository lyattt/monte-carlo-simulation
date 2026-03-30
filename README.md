# cmda3634-pr1 Concurrent Monte Carlo Simulation in Go


## Overview


This project implements a concurrent Monte Carlo simulation engine in Go. It simulates the growth of an investment portfolio over 252 trading days using randomly generated daily returns from a normal distribution. The program uses goroutines and channels to run many independent simulation trials in parallel and then aggregates the results to compute summary statistics.


## Go source Code Implementing the Engine


The engine is implemented in main.go and includes:

-   A function to simulate one trial of the stochastic process
-   Worker goroutines that execute trials in parallel
-   A concurrent execution function using channels to collect results
-   Statistical calculations (mean, standard deviation, and 95% confidence interval)
-   A command-line interface for user input

## How to Run the program 


Run the program using:

go run main.go -trials=10000 -workers=4

Optional (set a seed for reproducibility):

go run main.go -trials=10000 -workers=4 -seed=123

Arguments:

-trials: total number of simulation trials

-workers: number of goroutines to spawn

-seed: optional random seed


### Example Output 


Monte Carlo Simulation Results
Trials: 10000
Workers: 4
Time to run: 0.5 seconds
Mean outcome: 112.34
Std dev: 5.67
95% CI: [111.90, 112.78]


## Example Code for Benchmark Problem 


func simulateOnetRIAL() float64 {
   value := 100.0
   mu := 0.001
   sigma := 0.02

   for i := 0; i < 252; i++ {
      r := mu + sigma*rand.NormFloat6()
      value = value * (1 + r)
   }

   return value
}

## Extensions Attempted 


No additional extensions were implemented 


## AI Usage 


ChatGPT was used as a learning and debugging tool to help implement the concurrent execution portion of the project. It provided guidance on using goroutines and channels, structuring worker functions, and handling common errors during development. All code was reviewed, tested, and integrated manually.


## Benefits of AI Use


-  Helped understand concurrency concepts in Go
-  Assisted with debugging and resolving errors
-  Provided guidance on structuring parallel execution


## Limitations of AI Use 


-  Required manual verification to ensure correctness
-  Some suggestions needed adjustment to fit project requirements
-  Did not replace the need for testing and understanding the code


## AI Logs 


1. Experiment Function
   - [AI Log](experiment_ai_log.md)
2. Concurrent Execution
   - [AI log](concurrent_ai_log.md)
3. Aggregation
   - [AI log](aggregation_ai_log.md)
4. UserInterface(CLI)
   -[AI log](experiment_ai_log.md)
