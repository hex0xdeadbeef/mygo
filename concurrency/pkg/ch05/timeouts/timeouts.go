package timeouts

/*
It's recommended to place timeouts on all of our concurrent operations to guaranteeour system won't deadlock. The timeout period doesn't have to be close to the actual time it takes to perform
your concurrent operation. The timeout period's purpose is only to prevent deadlock, and so it only needs to be short enough that a deadlocked system will unblock in a reasonable amount of time
for our case.
*/