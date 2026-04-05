## Prompt 
How do I generate random values from a normal distribution in Go?

## Response
Summary of Chats answer: I used rand.NormFloat64() and scaled it with mu + sigma * rand.NormFloat64() to match the given mean and standard deviation.


## Prompt
 Why do we multiply by (1 + r) instead of adding returns?


 ## Response 
Summary of chats answer: This represents compounding returns. Each day’s return is applied as a percentage of the current value.