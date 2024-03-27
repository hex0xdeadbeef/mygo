package errpropagation

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime/debug"
)

/*
The information errors must have:

1) What happened.

For example "disk full", "socket closed", "credentials expired" and others.

2) When and where it occured.
Errors should always contain a complete stack trace starting with how the call was initiated and ending with where the error was instantiated. The stack
trace shouldn't be contained in the error message, but should be easily accessible when handling the error up the stack.

Further, the error should contain information regarding the context it's running within. For example: in a distributed system, it should have some way of identifying what machine the
error occured on, Later, when trying to understand what happened in your system, this information will be invaluable.

In addition, the error should contain the time on the machine the error was instantiated on, in UTC.

3) A friendly-facing message.

The error should contain abbreviated and relevant information from the previous two points. A friendly message is human-centric, gives some indication of
whether the issue is transitory, and should be about one line of text.

4) How the user can get more information.
Errors that are presented to users should provide an ID that can be cross-referenced to a corresponding log that displays the full information of the error:
	1) time the error occured (no the time error was logged)
	2) the stack-trace
	3) the hash of the stack trace to aid in aggregating like issues in bug trackers.
*/

/*
It's possible to place all errors into one of two categories:
	1) Bugs
	2) Known-edge cases (e.g. broken network connections, failed disk writes, etc.)
*/

/*

CLI Component -> Intermediary Component -> Low Level Component

Each component.
All incoming errors must be wrapped in a well-formed error for the component our code is within.

package intermediarycomp

import ".../lowlevel"

func PostReport(id string) error {
	result, err := lowlevel.DoWork()
	if err != nil {
		if _, ok := err.(lowlevel.Error); ok {
			err = WrapErr(err, "cannot post report with id %q", id)
		}
		return err
	}
	...
}
*/

/*
NOTE: That is only necessary to wrap errors in this fashion at your own module boundaries - public functions/methods or when your code can add valuable context.

Error correctness becomes an emergent propery of our system. We also concede perfection from the start by explicitly handling malformed errors, and by doing so we have given ourselves
a framework to take mistakes and correct them over time.

*/

/*
Remember that in either case, with well-or malformed errors, we will have included a log ID in the message to give the user something to refer back to should the user want more
information
*/

type MyError struct {
	Inner      error
	Message    string
	StackTrace string
	Misc       map[string]interface{}
}

func wrapError(err error, messagef string, msgArgs ...interface{}) MyError {
	return MyError{
		Inner:      err,
		Message:    fmt.Sprintf(messagef, msgArgs...), // messagef is a message format
		StackTrace: string(debug.Stack()),
		Misc:       make(map[string]interface{}), // This is where we might store the concurrent ID, a hash of the stack trace, or other contextual information that might help in diagoning
		// the error.
	}
}

func (me MyError) Error() string {
	return me.Message
}

// "lowlevel" module

type LowLevelErr struct {
	error
}

func isGloballyExec(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, LowLevelErr{(wrapError(err, err.Error()))}
	}

	return info.Mode().Perm()&0100 == 0100, nil
}

// "intermediate" module
type IntermediateErr struct {
	error
}

func runJob(id string) error {
	const (
		jobBinPath = "/bad/job/binary"
	)

	isExecutable, err := isGloballyExec(jobBinPath)
	if err != nil {
		return IntermediateErr{wrapError(err, "cannot run job %q: requisite binaries not available", id)}
	} else if !isExecutable {
		return wrapError(nil, "job binary is not executable")
	}

	return exec.Command(jobBinPath, "--id="+id).Run()
}

// main
func handleError(key int, err error, message string) {
	log.SetPrefix(fmt.Sprintf("[logID: %v] ", key))
	log.Printf("%#v", err)
	log.Printf("[%v] %v", key, message)
}

func Using() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	err := runJob("1")
	if err != nil {
		msg := "There was an unexpected issue: please report this as a bug."
		if _, ok := err.(IntermediateErr); ok {
			msg = err.Error()
		}
		handleError(1, err, msg)
	}
}
