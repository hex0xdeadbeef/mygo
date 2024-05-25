package main

// Анатомия каналов в Go - Хабр: https://habr.com/ru/articles/490336/
/*
CHANNEL CREATION ---------------------------------------------------------------------------------------------------------------------------------
1. The channel in Go is a pointer.


CHANNELS IN PRACTICE ---------------------------------------------------------------------------------------------------------------------------------
1. The channel in Go is a pointer.


DEADLOCK ---------------------------------------------------------------------------------------------------------------------------------
1. If we try to read from a channel, but this channel is empty, the scheduler will block reading goroutine and unlocks another assuming that another goroutine pass the data on this
channel. Similarly this happens in a case of sending data with the fullfilled channel: the scheduler will block the passing goroutine, until another one reads the data from the channel.


BUFFERED CHANNELS ---------------------------------------------------------------------------------------------------------------------------------
1. When the size of buffer is greater than 0, goroutine doesn't block until the buffer is full.

2. When buffer is full any vals are being sent through are added throwing the previous val available to be read (where a G isn't locked).
	1) The operation of reading from a BUFFERED channel is gready. It will end after full draining. It means the G will be reading from the channel without a lock until the buffer gots
	emptied.

CHANNEL WITH CHANNEL TYPE ---------------------------------------------------------------------------------------------------------------------------------


DEFAULT OPERATOR ---------------------------------------------------------------------------------------------------------------------------------
1. When a select has select statement, main() goroutine will be blocked and all others goroutines will be planned (one at a time) used in select statement.

2. Since select isn't blocked when it has the default statement, the planner doesn't launch all the available goroutines, but it can be done by calling time.Sleep(). After this all the
goroutines will be launched and when the managing will be returned to the main() goroutine, the channels will probably have the vals.

3. We can use default to prevent deadlock.


NIL CHANS ---------------------------------------------------------------------------------------------------------------------------------
1. The default value for a chan is nil, due to the fact that the chan is a pointer

2. Sending on a nil channel blocks the caller.


EMPTY SELECT  ---------------------------------------------------------------------------------------------------------------------------------
1. empty select statement results in deadlock


WAIT GROUP ---------------------------------------------------------------------------------------------------------------------------------
1. WaitGroup must be passed as a pointer to a structure.

2. Method Wait() is used to lock a holder of WaitGroup instance until the counter is decremented to 0

3. Internally the Done() method calls Add(-1)


WORKERS POOL ---------------------------------------------------------------------------------------------------------------------------------
1. Worker Pool is a set of goroutines that are working on a specific task. 	
*/
