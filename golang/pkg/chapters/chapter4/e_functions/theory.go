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
*/
