package fanoutin

/*
One of the interesting properties of pipelines is the ability they give us to operate on the stream of data using a combination of separate, often reorderable stages. We can even reuse
stages of the pipeline multiple times.
*/

/*
Fan-out is a term to describe the process of starting multiple goroutines to handle input from the pipeline. Input from the pipeline -> many goroutines
*/

/*
Fan-in is a term to describe the process of combining multiple results into one channel. Many results -> one channel
*/

/*
Conitions so that a stage suits for utilizing fan-out:
	1) The stage doesn't rely on values that the stage had calculated before.
	2) It takes a long time to run.
*/

/*
The property of order-independence is important because we have no guarantee in what order concurrent copies of our stage will run, nor in what order they will return.
*/