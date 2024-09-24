package main

/*
	USE METHODS ALLOWING TO PASS BYTE SLICES
1. It's better to use methods allowing to pass byte slice: these methods usually give more flexible control over distribution.
2. A good example of comparison: time.Format and time.AppendFormat.
	- The first one returns a string (under the hood it allocates a byte slice and calls the function time.AppendFormat
	over it).
	- The second one takes a byte slice, writes formatted time representation and returns expended byte slice.
3. Why does it rise the performance?
	- By using it we can pass byte slices that we've been given by sync.Pool instead of allocating new buffers
	- We can allocate a buffer of the length we need prematurely
*/