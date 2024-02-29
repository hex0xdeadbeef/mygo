package chapter11

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
11. TESTING
1. Testing (automated testing) is the practise of writing small programs that check that the code under test (the production code) behaves as expected for certain inputs, which are
usually either carefully chosen to exercise certain features or randomized to ensure broad coverage.

2. The task of testing occupies all programmers some of the time.

3. We have to be careful of boundary conditions, think about data structures, and reason about what results a computation should produce from suitable inputs.
------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
11.1 THE "go test" TOOL
1. The "go test" subcommand is a test driver for Go packages that are organized according to certain conventions.

2. In a package directory, files whose names ends with "_test.go" aren't part of the package ordinarily built by "go build" but a part of it when it built by "go test"

3. Within "*_test.go" three kinds of functions are treated specially:
	1) tests
	A test function which is a function whose name begins with "Test" exercises some program logic for correct behavior. "go test" calls the test function and reports the result, which is
	either "PASS" or "FAIL".

	2) benchmarks
	A benchmark function has a name beginning with "Benchmark" and measures the performance of some operation. "go test" reports the mean execution time of the operation.

	3) examples
	An example function whose name starts with "Example" provides machine-checked documentation.

The "go test" tool scans the "*_test.go" files for these special functions, generates a temporary main package that calls them all in proper way, builds and runs it, reports the results
and then cleans up.
------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
11.2 TEST FUNCTIONS
1. Each test file must import "testing" package.

2. The signature of test functions:
	func TestName(t *testing.T) {
		...
	}
		1) "Name" is an optional name and must begin with a capital letter.

3. "go test" (or "go build") command with no package arguments operates on the package in the current directory.

4. It's good practice to write the test first and observe that it triggers the same failure described by the user's bug report. Only then can we be confident that whatever fix we
come up with addresses the right problem.

5.Flags:
	1) The "-v" flag prints the name and execution time of each test in the package.

	2) The "-run" with regular expression argument causes "go test" run only those tests whose function name matches the pattern, so we can point which test functions "go test" must
	run. For example:
		go test -v -run="TestPalindrome|...|"

6. We can create place-in struct with all necessary fields consequently filling it to make a test suit. For example:
	tests = []struct {
		input    string
		expected bool
	}{
		{"aba", true},
		...
	}
The "table-driven" testing style is very common in go. It's straightforward to add new table entries as needed, and since the assertion logic isn't duplicated we can invest more effort
in producing a good error message.

7. Tests are independent of each other and that's the reason why the output of a failing test doesn't include the entire stack trace at the moment of the call to "t.Errorf". Also
"t.Errorf" doesn't cause a panic or stop the execution of the test, unlike assertion failures in many test frameworks for other languages. So if an early entry in the table causes test
to fail, later table entries will still be checked, and thus we may learn about multiple failures during a single run.

8. If we do need to stop all the tests, we use "t.Fatal" or "t.Fatalf". These must be called from the same goroutine as the "Test" function, not from another one created during the test.

9. Test failure messages are usually of the form: f(x) = y, want z, where:
		1) f(x) explains the attempted operation and its input
		2) y is the actual result
		3) z is the expected result

	1) Where convenitent Go syntax is used for the f(x) part.
	2) We should avoid boilerplate and redundant information.
	3) When testing a boolean function such as "IsPalindrome()", omit the "want z" part since it adds no information.
	4) if x,y,z is lengthy, print a concise summary of the relevant parts instead.

------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
11.2.1 RANDOMIZED TESTING

1. "Randomized testing" explores a broader range of inputs by constructing inputs at random.

2. How do we know what output to expect from our function, given a random input?
	1) To write an alternative inplementation of the function that uses a less efficient but simpler and clearer algorithm, and check that both implemetations give the same result.
	2) The second is to create input values according to a pattern so that we know what output to expect.

3. Since randomized tests are nondeterministic (can produce different results with the same arguments), it's critical that the log of the failing test record sufficient information
to reproduce the failure. For more complex functions that accept more complex inputs, it may be simpler to log the seed of pseudo-random number generator (as we do in
"TestRandomPalindrome()") than to dump the entire input data structure. Armed with that seed value, we can easily modify the test to replay the failure deterministically.

4. By using the current time as a source of randomness, the test will explore novel inputs each time it's run, over the entire course of the lifetime. This is especially valuable if
your project uses an automated system to run all its tests periodically.
------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
11.2.2 TESTING A COMMAND (AN EXECUTABLE FILE)

1. We can use "go test" for testing commands as well. A package named "main" ordinarily produces an executable program, but it can be imported as a library too.

2. The test code is in the same package as the production code. Although the package name is main and it defines a "main" fucntion, during testing this package acts a library that
exposes the function "TestEcho()" to the test driver. Its main function is ignored.

ERRORS
1.It's important that code being tested not call "log.Fatal()" or "os.Exit()", since these will stop the process in it tracks; calling these fucntions should be regarded as the
exclusive right of "main"

2. If something totally unexpected happens and a function panics, the test driver will recover, though the test will of course be considered a failure.

3. Expected errors such as those resulting from bad user input, missing files, or improper configuration should be reported by returning a non-nil "error" value.
------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
11.2.3 WHITE-BOX TESTING

1. One way of categorizing tests is by the level of knowledge they require of the internal workings of the package under test.

	1. A "black-box" test assumes nothing about the package other than what's exposed by its API and specified by its documentation. The package's internals are opaque.
		1) "black-box" tests are usually more robust, needing fewer updates as the software evolves.
		2) "black-box" tests also help the test author emphasize with the client of the package and can reveal flaws in the API design.

	2. A "white-box" (The name white box is traditional, but "clear box" would be more accurate) test has privileged access to the internal functions and data structures and can make
	observations and changes that an ordinary client cannot. For example: a "white-box" test can check that the invariants of the package's data types are maintained after every operations.
		1) "white-box" tests can provide more detailed coverage of the trickier parts of the implemetation.

2. While developing "TestEcho()", we modified "echo" function to use the package-level "out" when writing its output, so that the test could replace the standart output with an
alternative implemetation that records the data for later inspection. With the same approach, we can replace other parts of production code with easy-to-test "fake" implemetations.
The advantages of "fake" implemetations are:
	1) They can be simpler to configure
	2) They are more predictable, reliable, and easier to observe.
	3) They can also avoid undesirable side effects such as updating a production code database or charging a credit card.

3. We'd like to test "storage" but we don't want the test to send out real email. So we move the email logic to its own function and store that function in an unexported package level
variable, "notifyUser()". Now we can write a test that substitutes a simple fake notification mechanism instead of sending real email.
	1) There's one problem. After "TestCheckQuouta()" function has returned, "CheckQuota()" no longer works it should because it's still using the test's fake implemetation of
	"notifyUser()" (There's always the risk of this kind when updating global variables). We must modify the test to restore the previous value so that subsequent tests observe no
	effect, and we must do this on all execution paths, including test failures and panics. This naturally suggests "defer"
This pattern can be used to temporarily save and restore all kinds of global variables including:
	1) command-line Flags
	2) debugging options
	3) performance parameters
to install and remove hooks that cause the production code to call some test code when something interesting happens and to coax the production code into rare but important states,
such as:
	1) timeouts
	2) errors
	3) specific interleavings of concurrent activities

4. Using global variables in this way is safe only because "go test" doesn't normally run multiple tests concurrently.
------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
11.2.4 EXTERNAL TEST PACKAGES
1. Consider the "net/url" and "net/http" packages. "net/http" depends on the lower-level "net/url". However, one of the tests in "net/url" is an example demostrating the interaction
between URLs and the HTTP client library. In other words there's a cyclic dependencies. A test of the lower-level package imports the higher-level package.
	1) We solve this problem declaring the test function in an "external test package", in a file in the "net/url" directory whose package declaration reads "package url_test".
	The extra suffix "_test" is a signal to "go test" that it should build an additional package containing just these files and run its tests. it may be helpful to think of this
	external package as if it had the import path "net/url_test", but it cannot be imported under this or any other name.

	2) Because external tests live in a separate package, they may import helper packages that also depend on the package being tested. An in-package test cannot do this. In terms
	of the design layers, the external test package is logically higher up than both of the packages it depends upon.

	3) By avoiding import cycles, external test packages allow tests, especially, "integration tests" (which test the interaction of several components), to import other packages
	freely, exactly as an application would.

2. We can use use the "go list" tool to summarize which Go source files in a package are in-package tests or external tests
	1) -f={{.GoFiles}} is the list of that contain production code. These are the files that "go build" will include in our application.
	2) -f={{.TestGoFiles}} is the list of files that also belong to the "fmt" package, but these files whose names all end in "_test.go" are included only when building tests.
	3) -f={{.XTestGoFiles}} is the list of files that constitute the external test package

3. Sometimes an external test package may need privileged access to the internals of the package under test, if for example a white-box test must live in a separate package to
avoid an import cycle. In such cases, we use a trick:
	1) We add declarations to an in-package "_test.go" file to expose the necessary internals to the externals test. This file thus offers the test a "back door" to the package. If
	the source file exists only for this purpose and contains no test itself, it's often called "export_test.go"
For example:
The implemetation of the "fmt" package needs the functionality of "unicode.IsSpace()" as part of "fmt.Scanf()". To avoid creating and undesirable dependency, "fmt" doesn't import
the "unicode" package and its large tables of data. Instead, it contains a simpler implementation, which it calls "isSpace". To insure that the behaviors of "fmt.isSpace" and
"unicode.IsSpace" don't drift apart "fmt" prudently contains a test. There's an external test in package "fmt_test" and it cannot access "isSpace()" directly, so "fmt" opens a "back
door" to it by declaring an exported variable that holds the internal "isSpace" function. This is the entriety of the "fmt" package's "export_test.go" file:
	package fmt
	var IsSpace = isSpace
This test defines no tests; it just declares the exported symbol "fmt.IsSpace()" for use by the external test. This trick can also be used whenever an external test needs to use
some of the techniques of white-box testing.
------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
11.2.5 WRITING EFFECTIVE TESTS

1. Go's attitude to testing stands in stark contrust. It expects test authors to do most work themselves, defining functions to avoid repetition, just as they would for ordinary
programs. The process of testing is not one of rote form form filling; a test has a user interface too, albeit one whose only users are also ist maintainers.

2. A good test doesn't explode on failure but prints a clear and succint description of the symptom of the problem, and perhaps other relevant fact about the context. Ideally, the
maintainer shouldn't need to read the source code to decipher a test failure.

3. A good test shouldn't give up after one failure but should try to report several errors in a single run, since the pattern of failures may itself be revealing.

4. Using "TestSplitWithAssertion()" we forfeit the opportunity to provide meaningful context. Using a full context test function that reports we emphasize the significance of the
result. Full context function identifies the actual value and the expectation and it continues to execute even if this assertion should fail.

5. The key to a good test is to start by implementing the concrete behavior that you want and only then use function to simplify the code and eliminate repetition.
------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
11.2.6 AVOIDING BRITTLE TESTS

1. An application that often fails when it encounters new but valid inputs is called "buggy"

2. A test that spuriously fails when a sound change was made to the program is called "brittle". The most brittle tests, which fail for almost any change to the production code are
sometimes called "change detector" or "status quo" tests, and the time spent dealing with them can quickly deplete any benefit they once seemed to provide

3. Just a buggy program frustrates its users, a brittle one exasperates

4. The easiest way to avoid brittle tests is to check only the properties you care about. Test your program's simpler and more stable interfaces to its internal functions. Don't
check for exact string matches, but look for relevant substrings that will remain unchanged as the program evolves. It's often worth writing a substantial function to distill a
complex output down to its essence so that assertions will be reliable.
------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
11.3 COVERAGE
1. No quantity of tests can even prove a package free of bugs. At best, they increased our confidence that the package works well in a wide range of important scenarious.

2. The degree to which a test suite exercises the package under test is called "test's coverage". Coverage cannot be quantified directly - the dynamics of all but the most trivial
programs are beyond precise measurement - but there are heuristics that can help us direct our testing efforts to where they are more likely to be useful.

3. "Statement coverage" is the simplest and most widely used of these heuristics. The statement coverage of a test suite is the fraction of source statements that are executed at least
once during the test.

4. In this section we'll use Go's "cover" tool, which is integrated into "go test".

5. The command "go tool cover" displays the usage message of the coverage tool. The "go tool command" runs one of the executables from the Go toolchain. These programs live in the directory $GOROOT/pkg/tool/${GOOS}_${GOARCH}. Thanks to "go buid", we rarely need to invoke them directly. In sum we can run the test with the "-coverprofile" flag:
	go test -run=TestCoverage -coverprofile=c.out
		1) When c.out - filename in which the data will be written out

This flag enables the collection of coverage data by instrumenting the production code. That is, it modifies a copy of the source code so that before each block of statements is executed,
a boolean variable is set, with ove variable per block. Just before the modified program exits, it writes the value of each variable to the specified log file "c.out" and prints a summary
of the fraction of statements that were executed. If we need only a percent of coverage we should use:
	go test -cover

6. If "go test" is run with the "covermode=count" flag, the instumentation for each block increments a counter instead of setting a boolean. The resulting log of execution counts of
each block enables quantitative comparisons between "hotter" blocks, which are more frequently executed, and "colder" ones.

7. Having gathered the data, we run the cover tool, which processes the log, generates an HTML report, and opens a new browser window:
	go tool cover -html=c.out
		1) When c.out - filename in which the result data is located

8. Achieving 100% statement coverage sounds like a noble goal, but it's not usually feasible in practice, nor is it likely to be a good use of effort.

9. Testing is fundamentally a pragmatic endeavor, a trade-off between the cost of writing tests and the cost of failures that could have been prevented by tests.

10. Coverage tools can help identify the weakest spots, but devising good test cases demands the same rigorous thinking as programming in general.

------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/

/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
11.4 BENCHMARK FUNCTIONS
1. Benchmarking is the practise of measuring the performance of a program on a fixed workload.

2. Benchmark fuction has the "Benchmark" prefix and "*testing.B" parameter that provides most of the same methods as a "*testing.T". It also exposes an integer field "N", which specifies
the number of times to perform the operation being measured.

3. Benchmarks are run with the command "go test -bench=.". By default o benchmarks are run. The argument to the "-bench" flag selects which benchmarks to run. It's a regular expression
matching the names of "BenchmarkXxx" functions, with a default value that matches none of them.
	1) The "." pattern causes it to match all benchmarks

The outut is:
	BenchmarkIsPalindtome-12         7797498               146.8 ns/op           160 B/op          2 allocs/op

	1) BenchmarkIsPalindtome-12 - is a value of GOMAXPROCS (the number of CPUs)
	2) 7797498 - the nubmer of operations. Since the benchmark runner initialyy has no idea how long the oreration takes, it makes some initial measurements using small values of N, and
	then extrapolates to a value large enough for a stable timing measurement to be made.
	3) 146.8 ns/op the time elapsed for each operation
	4) 
	5) 2 allocs/op - allocation statistic

4. The reason the loop is implemented by the benchmark function, and not by the calling code in the test driver is so that the benchmark function has the opportunity to execute any 
necessary one-time setup code outside the loop without this adding to the measured time of each iteration.
	1) If this setup code is still perturbing the results, "the testing.B" parameters provides methods to stop, resume and reset the timer, but these are rarely needed.

5. The "-benchmem" command-line flag will include memory allocation statistics in its report. 

6. Consolidating the allocations in a single call to "make()" eliminates 75% of the allocations and halves the quantity of allocated memory.

7. In many settings the interesting performance questions are about the relative timings of two different operations. 
1) For example:
	If a function takes 1ms to process 1000 elements, how long will it take to process 10,000 or 1,000,000?

Such comparisons reveal the asymptotic growth of the running time of the function. 
2) Another example:
	What's the best size for an I/O buffer? Benchmarks of application throughput over a range of sizes can help us choose the smallest buffer that delivers satisfactory performance.

3) A third example
	Which algorithm performs best for a given job? Benchmarks that evaluate two different algorithms on the same data can often show the strengths and weakenesses of each one on important
	and representative workloads.

7. Comparative benchmarks are just regular code. They typically take the form of a single parameterized function, called from several "BenchmarkXxx" functions with different values, like
this ones:
	func benchmark(b *testing.B, size int) { ... }
	Benchmark10(b *testing.B) { benchmark(b, 10) }
	...
	Benchmark1000(b *testing.B) { benchmark(b, 1000) }
Resist the temptation to use the parameter "b.N" as the input size. Unless you interpret it as an iteration count for a fixed-size input, the results of your benchmark will be 
meaningless.

8. As the program evolves, or its imput grows, or it is deployed on new operating systems or processors with different characteristics, we can reuse the benchmarks to revisit design 
decisions.
------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/


/*------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
11.5 PROFILING

1. Donald Knuth's aphorism about premature optimisation: "Structured Programming with go to Statements". The meaning is:
	There's no doubt that the grail (the blood of the Jesus's children) of efficiency leads to abuse. Programmers waste enormous amounts of time thinking about, or worrying the speed
	of noncritical parts of their programs, and these attempts at efficiency actually have a strong negative impact when debugging and maintenance are considered. WE SHOULD FORGET 
	ABOUT SMALL EFFICIENCIES, SAY ABOUT 97% OF THE TIME: PREMATURE OPTIMIZATION IS THE ROOT OF EVIL.

	Yet we should not pass up our opportunities in that critical 3%. A good programmer won't be lulled into complacency by such reasoning, he'll be wise to look carefully at the
	critical code, but only after that code has been identified, It's often a mistake to make priori judgments about what parts of a program are really critical, since the universal
	experience of programmers who have been using measurement tool has been that their intuitive guesses fail.

2. When we wish to look carefully st the speed of our programs, the best technique for identifying the critical code is "profiling". Profiling is an automated approach to performanse
measurement based on sampling a number of profile events during execution, then extrapolating from them during post-processing step; the resulting statistical summary called "profile".

3. Go supports many kinds of profiling, each concerned with a different aspect of performance, but all of them involve recording a sequence of interest, each of which has an accompanying
stack trace. The "go tool" has buit in support for several kinds of profiling.
	1) CPU PROFILE identifies the functions whose execution requires the most CPU time. The currently running thread on each CPU is interrupted periodically by the operating system every
	few milliseconds, with each interruption recording one profile event before normal execution resumes.
	2) HEAP PROFILE identifies the statements responsible for allocating the most memory. The profiling library samples calls to internal memory allocation routines so that on average, one
	profile event is recorded per 512KB of allocated memory.
	3) BLOCKING PROFILE identifies the operations responsible for blocking goroutines the longest, such as system calls, channel sends and receives, and acquisitions of locks. The profiling
	library records an event every time a goroutine is blocked by one of these operations.

4. So that we can use profiling we should enable some of these flags. But we must be careful when using more than one flag at a time, however: the machinery for gathering one kind of profile
may skew the results of other.
	1) go test -cpuprofile=cpu.out
	2) go test -blockprofile=block.out
	3) go test -memprofile=mem.out

5. It's easy to add profiling support to non-test programs too, though the details of how we do that vary between short-lived command tools and long-running server applications. Profiling
is especially useful in long-running applications, so the Go runtime's profiling features can be enabled under programmer control using the "runtime" API.

6. Once we've gathered a profile, we need to analyze it using the "pprof" This is a standart part of the Go distribution, but since it's not an everyday tool, it's accessed indirectly using
"go tool pprof". It has dozens of features an options, but basic use requires only two arguments:
	1) The executable that produced the profile
	2) Profile log. To make profiling efficient and to save space, the log doesn't include function names. Instead functions are identified by their addresses. This means that "pprof" needs
	the executable in order to make sense of the log. Although "go test" usually discards the test executable once the test is complete, when profiling is enabled it saves the executable 
	as "packageName.test".

7. For example:
	1) go test -run=NONE -bench=ClientServerParallelTLS64 -cpuprofile=cpu.log net/http
	2) go test -run=NONE -bench=ClientSeverParallelTLS64 -cpuprofile=cpu.log net/http 
------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/