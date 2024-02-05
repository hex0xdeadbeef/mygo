package e_functions

/*
FUNCTIONS-------------------------------------------------------------------------------------------------------------------------------
1. Function results may be also named. In this case, each name declares a local variable initialized to the zero value
of its type.

2. We can omit type declaration for each parameter substitution it with the single right type declaration

3. A blank parameter "_" may be used to emphasize that a parameter is unused.

4. The type of functions is also called as signature. Two functions have the same type if they have the same sequence
of parameter types and the same sequence of result types. The names of parameters and results don't affect the type,
nor does whether or not they were declared using factored form.

5. Parameters are local variables within the body of the function, with their initial values set to the arguments sup
plied by the caller.

6. Arguments are passed by value, so the function receives a copy of each argument. Modifications to the copy don't
affect the caller.
! If an argument has some arguments contain some kind of reference, like a pointer, slice, map, function or channel,
then the caller may be affected by any modifications the function makes to variables indirectly referred to by the ar
gument

7. If a function hasn't got a body, it means that it's implemented in a language other than Go:
	func Sin(x float64) float64 - implemented in assembly language

8. Call stack size is running from 64 KB to 2MB.
! Go uses variable-size stacks that start small and grow as needed up to a limit on the order of a gigabyte.

9. A function can return more than one result.
	1) The result of calling a multiple-valued function is a tuple and it must be explicitly assigned to variables.
	! To ignore some of them we use "_" !
	2) If we don't interact with multiple named parameters they values will have the zero values of its types.
	3) In a function with named results return operands may be omitted.
	Due to this fact the function can return some zero values of its types. So that we supply our code with simplicity
	we should avoid this pattern.


ERROR HADNLING-------------------------------------------------------------------------------------------------------------------------------
1. A function for which failure is an expected behavior returns an addition result, conventionally the last one.
	1) If the failure has only one possible cause, the result is a boolean, called "ok"

2. The built-in type error is an interface type.
	1) Error can be nil or non-nil
		1) nil implies success

		2) non-nil implies failure and this type of errors must be describe with a message is sent to a caller.
		We can obtain it by calling "Error()" method, fmt.Println(err) or fmt.Printf("%v", err).

3. Usually when a function return a non-nil error other results are undefined and should be ignored.

4. If a function returns error and partial processed data, this function should be documented properly. The unprocessed data
should be handled by a caller.

5. Error handling flavors are:
	1) Propagating the error, so that a failure in a subroutine becomes a failure of the calling routine. If the subroutine
	with failure doesn't provide sufficient information, we need to add descriptive information to an error and return it.
	Also we must create a chain of errors from "the depths to the surface".
 		! DUE TO CHAINED ERRORS HANDLING ERRORS SHOULDN'T BE CAPITALIZED AND NEWLINES SHOULD BE AVOIDED !
		1) In general the call f(x) is responsible for reporting the attempted operation f and the argument value x as they
		relate to the context of the error.

		2) The caller is responsible for adding further information that it has but the call f(x) doesn't

	2) For errors that represent transient or unpredictable problems, it makes sense to retry the failed operation possibly
	with:
		1) A delay between tries

		2) With a limit on the number of attempts

		3) Time spent trying before giving up entirely

	Finally if progress is impossible, the caller can print the error and stop the program gracefully. But this course of
	action should generally be reserved for the main package of a program (in particular int the main function.

	! LIBRARY FUNCTIONS SHOULD USUALLY PROPAGATE ERRORS TO THE CALLER, UNLESS THE ERROR IS A SIGN OF AN INTERNAL
	INCONSISTENCY !

6. As a way to issue the error after non-nil result we can use "fmt.Fprintf()"/"log.Fatalf()" which constructs formatted output,
by default it prefixes the time and date to the error message.
	1) The default "fmt.Fprintf()" is helpful for a long-running server.
	2) The second "log.Fatalf()" is helpful for an interactive tool

7. In order to get a more attractive output, we can set the prefix used by the "log" package to the name of the command, and suppress the
display of the date and time.
		E.G:
		1) log.SetPrefix("wait: ")

		2) log.SetFlags(0)

8. In some cases it's sufficient just to log the error and then continue, perhaps with reduced functionality.

9. In rare cases we can safely ignore an error entirely.
! DISCARDING THE ERRORS MUST BE DELIBERATE AND ITS INTENTION MUST BE DOCUMENTED MEANINGFUL !

10. Functions tend to exhibit a common structure, with a series of initial checks to reject errors, followed by the function at the end, minimal
ly indented.

EOF--------------------------------------------------------------------------------------------------------------------------------------------
1. If the caller repeatedly tries to read fixed-size chunks until the file is exhausted, the caller must respond differently to an end-of-file
condition than it does to all other errors. To take this requirement there is an EOF error is provided by "io" package.

2. EOF has a fixed error message "EOF"

FUNCTION VALUES---------------------------------------------------------------------------------------------------------------------------------
1. Function values have types.

2. Functions may be assigned to variables or passed to or returned from functions.

3. Function value may be called like any other function.

4. The zero value of a function is nil
	1) The nil function call will result in a panic error
	2) We can compare functions to nil value

5. It's forbidden to assign to a function variable another function that has different type (It's the same to basic types)

6. Functions aren't comparable
	1) We can't compare functions
	2) We can't put functions into a map as keys

7. Functions let us to parameterize our functions over not just data, but behavior too. It allows us to separate portions
of logic.

8."%*s" verb is used to indent the output using "*" argument to set a symbols count and "%s" as a symbol to be indented with.

ANONYMOUS--------------------------------------------------------------------------------------------------------------------------------------------
1. Inner returnable anonymous functions can refer to the outer function variables.

2. Every function isn't just a code. The function has a state.

3. When we return a function from another one, the result function saves inner function state. So, with a call
of returned function we can change the state of parent's function.

4. When an anonymous fuction requires the recursion call to itself, we must declare it before assigning. In this
case trying to combine these two steps in a single one the compiler tries to detect the named function thoughout a
source file and due to the fact of impossibility of this action the compiler marks the statement as an error.

CAPTURING LOOPS VARIABLE-------------------------------------------------------------------------------------------------------------------------
1. When we capture a particular variable within another function, in fact we force the recepient function to
refer to the local memory pool. Due to this fact while capturing iterable variable we must assign this variable
to another new variable with shorthand in order to create another mempool and our function refers to it.

VARIADIC FUNCTIONS--------------------------------------------------------------------------------------------------------------------------------------------
1. To declare a variadic function, the type of the final parameter is preceeded by an ellipsis "..."

2. An variadic argument implies the slice of arguments.

3. In the case when arguments are passed directly the following actions happen.

Implicitly, the caller allocates an array, copies the arguments into it, and passes a slice of the entire array to the
function.
4. When arguments are already placed in slice: place an ellipsis after the final argument.

5. Although a variadic argument behaves like a slice within the body of the function the type of a function with this
parameter is different from the function with an explicit slice parameter.

DEFER--------------------------------------------------------------------------------------------------------------------------------------------
1. A defer statement is an ordinary function or method call prefixed by the keyword defer

2. Any number of calls may be deferred. They are executed in the reverse order of the order they
were deferred. LIFO (Last In First Out)

3. Defer is used with paired operations like open/close, connect/disconnect, lock/unlock to ensure
that all resourses are released in all cases, no matter how complex the control flow.

4. The right place for a defer statement that releases a resource is immediately after the resourse has been succes
sfully acquired.

5. The defer statement can also be used to pair "on entry" and "on exit" actions when debugging a complex function.

6. We should remember about extra parentheses after the declared deferred function.

7. If we defer the function that returns another one, the body of first function will be executed and
on the other hand the result will be deferred.

8. Deferred functions run AFTER return statement

9. Deferred anonymous function can observe the function's named results. It may be useful in functions
with many return statements.

10. Deferred functions can even change the values that the ecnlosing function returns to its caller.

11. Before deferring the values that interact with a loop variable we should fix a the loop variable in order to capture
a new one mempool.

12. While working with some files opening them into a loop scope and deferring their closing in the same place, we can run
out of all the file descriptors. To defeat this exhaustion we should enclose all work with file in a function.
	1) A file descriptor is a number a system uses to make an identification of an opened file.
	! THE DEFAULT NUMBER OF SIMULTANEOUS OPENED FILE DESCRIPTORS IS 1024 ON THE UNIX/LINUX SYSTEMS !

13. The io.Copy() function postpones any copying errors until the file is closed, so we should explicitly close the file and
track any errors.

PANIC--------------------------------------------------------------------------------------------------------------------------------------------
1. If Go's runtime detects mistakes as (out-of bounds/trying to access data with nil pointer) it panics.

2. Panic steps:
	1) The execution stops.
	2) All the deferred function calls in that goroutine are executed.
	3) The program crashes witg a log message. This value is usually an error message of some sort, and for
	each goroutine, a stack trace showing the stack of function calls that were active at the time of panic.

3. Not all panics come from the runtime.

4. We can explicitly cause a panic situation with the variadic "panic" funciton. It accepts any sort of arguments.

5. It's a good practise to assert that all the preconditions are held. If we cannot add anything more informa-
tive to the panic info, there's no point to do it.

6. In a robust program the following "expected" errors like incorrect input, misconfiguration, failing I/O
should be handled gracefully. They're best dealt with using error values.

7. When a panic occurs all deferred functions are run in reverse order, starting with those of the topmost function
on the stack and proceeding up to main.

8. It's possible for a function to recover from a panic so that it doesn't terminate the program.

9. For diagnostic purposes the runtime package lets the programmer dump the stack using the same machinery. By
deferring a call to printStack in main.

10. Go's panic mechanism runs the deferred functions before it unwinds the stack.

RECOVER--------------------------------------------------------------------------------------------------------------------------------------------
1. If "recover()" is called within the deferred function and the function contains the defer statement is panick
ing, recover ends the current state of panic and returns the panic value. The function that was panicking doesn't
continue where it left off but returns normally.

2. If recover is called in normal situation, it returns nil.

3. We might supply recover with a panic value.

4. We shouldn't attempt to recover from another package's panic.

5. Public APIs should report failures as errors.

6. We shouldn't recover from a panic that may pass through a function we don't control. E.g. the function from the
client/library 'cause we cannot ensure that its logic won't be broken and our data won't leak.

7. For some conditions there is no recovery. Running out of memory, for example, causes the Go runtime to terminate
the program with a fatal error.
*/
