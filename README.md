# cmda3634-pr1
Core Functionality
1. Experiment Function
   - [AI Log](experiment_ai_log.md)
   Question: How do I generate random values from a normal distribution in Go?
Summary of Chats answer: I used rand.NormFloat64() and scaled it with mu + sigma * rand.NormFloat64() to match the given mean and standard deviation.

• Question: Why do we multiply by (1 + r) instead of adding returns?
Summary of chats answer: This represents compounding returns. Each day’s return is applied as a percentage of the current value.
2. Concurrent Execution
   - [AI log](concurrent_ai_log.md)
3. Aggregation
   - [AI log](aggregation_ai_log.md)