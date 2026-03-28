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

func main() {
    result := simulateOneTrial()
    fmt.Println(result)
}