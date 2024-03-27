package goroutinehealing

/*
To heal goroutines, we'll use our heartbeat pattern to check up on the liveness of the goroutine we're monitoring.

The type of heartbeat will be determined by what we're trying to monitor, but if our goroutine can become livelocked, make sure that the heartbeat contains some kind of information
indicating that the goroutine is not only up, but doing useful work.
*/

/*
1. The logic that monitors a goroutine's health is a "steward". Stewards will also be responsible for restarting a ward's goroutine should it become unhealthy. To do so, it will need a
reference to a function that can start the goroutine.
2. The goroutine is monitored is called "ward".
*/
