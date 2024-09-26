package main

/*
	RUNTIME PACKAGE

1. Package runtime contains operations that interact with Go's runtime system, such as functions to control goroutines. It also includes the low-level type information used by the reflect package.

	ENVIRONMENT VARIABLES
1. The following environment variables ($name depending on the host operating system) control the run-time behavior of Go programs. The meanings and use may change from release to release.
2. The GOGC variables sets the initial garbage collection target percentage. A collection is triggered when the ration of freshly allocated data remaining after the previous collection reaches this percentage.

	The default percentage of the GOGC=100.

	Setting GOGC=off disables the garbage collector entirely.

	the runtime/debug.SetGCPercent allows changing thus percentage at run time.
2. The GOMEMLIMIT variable sets a soft memory limit for the runtime. This memory limit includes the
the Go heap and all other memory managed by the runtime and excludes exetrnal memory sources such as:
	- mappings of the binary itself
	- memory managed in other languages
	- memory held by the OS on behalf of the Go program

	GOMEMLIMIT is a numeric value in bytes with an optional unit suffix. The supported suffixes include:
		- B
		- KB
		- MiB
		- GiB
		- TiB
	These suffixes represent quantities of bytes as defined by the standart. That is, they're based on powers of two.

	The default value of this variable is math.MaxInt64, whic effectively disables the memory limit.

	runtime/debug.SetMemoryLimit allows changing this limit at run time.
3. The CODEBUG variable controls debugging variables within the runtime. It's a comma-separated list of name=val pairs setting these named variables.
	1. clobberfree
		setting clobberfree=1 causes the GC to clobber memory content of an object with bad content when it frees the object.
	2. cpu.*
		- cpu.all=off disables the use of all optional instruction set extensions
		- cpu.extension=off disables use of instructions from the instruction set extension suc as "sse41" or "avx" as listed in "internal/cpu" package. As an example: cpu.avx=off disables runtime detection and thereby
		use of AVX instructions.
	3. cgocheck
		- setting cgocheck=0 disables all checks for packages using cgo to incorrectly pass Go pointers to non-Go code.
		- setting cgocheck=1 (the default) enables relatively cheap checks that may miss some errors.

		A more complete, but slow, cgocheck mode can be enabled using GOEXPERIMENT (which requires a rebuild)
	4. disablehp
		- setting disablehp=1 on Linux disables transparent huge pages for the heap.
		It has no effect on other platforms. disablehp is meant for compatibility with versions of Go before 1.21, which stopped working around a Linux kernel default that can result in significant memory overuse. This
		setting will be removed in a future release, so operations should tweak their Linux configuration to suit their needs before then.
	5. dontfreezetheworld
		By default, the start of a fatal panic or throw "freeze the world", preempting all threads to stop all running goroutines, which makes it to traceback all goroutines, and keeps their state close to the point of
		panic.

		Setting dontfreezetheworld=1 disables this preemption, allowing goroutines to continue executing during panic processing. Note that goroutines that naturally enter the scheduler will stop.

		This can be useful when debugging the runtime scheduler, as freezetheworld perturbs scheduler state and thus may hide problems.
	6. efence
		- setting efence=1 causes the allocator to run in a mode where each object is allocated on a unique page and addresses are never recycled.
	7. gccheckmark
		- setting gccheckmark=1 enables verification of the GC's concurrent mark phase by performing a second mark pass while the world is stopped.
		If the second pass finds a reachable object that wasn't found by concurrent mark, the GC will panic.
	8. gcpacetrace
		- setting gcpacetrace=1 causes the GC to print information about the internal state of the concurrent pacer
	9. gcshrinkstackoff
		- setting gcshrinkstackoff=1 disables moving goroutines onto smaller stacks.
		In this mode, a goroutine's stack can only grow.
	10. gcstoptheworld
		- setting gcstoptheworld=1 disables concurrent GC making every garbage collection a STW event.
		- setting gcstoptheworld=2 also disables concurrent sweeping after the garbage collection finishes.
	11. gctrace
		- setting gctrace causes the GC to emit a single line to standart error at each collection, summarizing the amount of memory collected and the length of the pause. The format of this line is subject to change.
		Included in the explanation below is also the relevant runtime/metric for each field:

			// ...

		The phases are stop-the-world (STW) sweep termination, concurrent mark an scan, and STW mark termination.

		The CPU times for mark/scan are broken down in to assist time (GC performed in line with allocation), background GC time, and idle GC time.

		If the line ends with "(forced)", this GC was forced by a runtime.GC() call.
	12. harddecommit:
		- setting harddecommit=1 causes memory that is returned to the OS to also have protection removed on it.

		This is the only mode of operation on Windows but is helpful in debugging scavenger-related issues on other platforms.

		Currently, only supported on Linux.
	13. inittrace
		- setting inittrace=1 causes the runtime to emit a single line to standart error for each package with init work, summarizing the execution time and memory allocation.

		No information is printed for inits executed as part of plugin loading and for packages without both user defined and compiler generated init work.

		The format of this line is subject to chagne:

		// ...
	14. medvdontneed
		- setting medvdontneed=0 will use MADV_FREE instead of MADV_DONTNEED on Linux when returning memory to the kernel.

		This is more efficient, but means RSS numbers will drop only when the OS is under memory pressure.

		On the BSDs and Illumos/Solaris, setting madvdontneed=1 will use MADV_DONTNEED instead of MADV_FREE. This is less efficient, but causes RSS numbers to drop more quickly.
	15. memprofilerate
		- setting memprofilerate=X will update the value of runtime.MemProfileRate.

		When set to 0 memory profiling is disabled

		Refer to the descr. of MemProfileRate for the default value
	16. profstackdepth
		profstackdepth=128 (the default) will set the maximum stack depth used by all pprof profilers except for the CPU profiler to 128 frames.

		Stack traces that exceed this limit will be truncated to the limit starting from the leaf frame.

		- setting profstackdepth to any value above 1<<10 will silently default to 1<<10.
		Future versions of Go may remove this limitation and extend profstackdepth to apply to the CPU profiler and execution tracer.
	17. pagetrace
		- setting pagetrace=/path/to/file will write out a trace of page events that can be viewed, analyzed, and visualized using the x/debug/cmd/pagetrace tool.

		Build your program with GOEXPERIMENT=pagetrace to enable this functionality. Don't enable this functionality if your program is a setuid binary as it introduces a security risk in that scenario.

		Setting this option for some applications can produce large traces, so use it with care.
	18. panicnil
		- setting panicnil=1 disables the runtime error when calling panic with nil interface value or untyped nil.
	19. runtimecontentionstacks
		- setting runtimecontentionstack=1 enables inclusion of call stacks related to contention on runtime-internal locks in the "mutex" profile, subject to the MutexProfileFraction setting.

		When runtimecontentionstacks=0 contention on runtime-internal locks will report as "runtime._LostContendedRuntimeLock".

		When runtimecontentionstacks=1, the call stacks will correspond to the unlock call that released the lock. But instead of the value corresponding to the amount of contention that call stack caused, it corresponds
		to the amount of time the caller of unlock had to wait in its original call to lock.

		A future release is expected to align those and remove this setting.
	20. invalidptr
		invalidptr=1 (the default) causes the GC and stack copier to crash the program if an invalid pointer value (for example 1) is found in a pointer-typed location.

		- Setting invalidptr=0 disables this check.

		This should only be used as a temporary workaround to diagnose buggy code.

		The real fix is not to store integers in pointer-typed locations.
	21. sbrk
		- setting sbrk=1 replaces the memory allocator and GC with a trivial allocatior that obtains memory from the OS and never reclaims any memory
	22. scavtrace
		- setting scavtrace=1 causes the runtime to emit a single line to standart error, roughly once per GC cycle, summarizing the amount of work done by the scavenger as well as the total amount of memory returned to
		the OS and an estimate of physical memory utilization. The format of this line is subject to change:

		// ...

		If the line ends with "(forced)", then scavenging was forced by a debug.FreeOSMemory()
	23. scheddetail
		- setting schedtrace=X and scheddetail causes the scheduler to emit detailed multiline info every X milliseconds, describing state of the scheduler, processors, threads and goroutines.
	24.	tracebackancestors
		- setting tracebackancestors=N extends tracebacks with the stacks at which goroutines were created, where N limits the number of ancestor goroutines to report.

		This also extends the information returned by runtime.Stack.

		Setting N to 0 will report no ancestry information.
	25. tracefpunwindoff
		- setting tracefpunwindoff=1 forces the execution tracer to use the runtime's default stack unwinder instead of frame pointer unwinding.

		This increases trace overhead, but could be helpful as a workaround or for debugging unexpected regressions caused by frame pointer unwinding.
	26. traceadvanceperiod
		traceadvanceperiod is the approximate period of nanoseconds between trace generations.

		Only applies if a program is built wotj GOEXPERIMENT=exectrace2.

		traceadvanceperiod is used primarily for testing and debugging the execution tracer.
	27. tracecheckstackownership
		- setting tracecheckstackownership=1 enables a debug check in the execution tracer to double-check stack ownership before taking a stacktrace.
	28. asyncpreemptoff
		setting asyncpreemptoff=1 disables signal-based asynchronous goroutine preemption.

		This makes some loops non-preemptible for long periods, which may delay GC and goroutine scheduling.

		This is useful for debugging GC issues because it also disables the conservative stack scanning used for asynchronously preempted goroutines.
4. GOMAXPROCS variable limits the number of OS's threads that can execute user-level Go code
	simultaneously.

	There is no limit to the number of threads that can be blocked in syscalls on behalf of Go code;
	those don't count against GOMAXPROCS limit.

	This package's GOMAXPROCS function queries and changes the limit.
5. The GORACE variable configures the race detector, for programs built using "-race".
6. The GOTRACEBACK variable controls the amount of output generated when a Go program fails due to
	an uncovereged panic or an unexpected runtime condition.

	By default, a failure prints a stack trace for the current goroutine, eliding functions
	internal to the run-time system, and then exits with exit code 2.

	The failure prints stack traces for all goroutines if there's no current goroutine or the failure is internal to the run-time.

	- GOTRACEBACK=none omits the goroutine stack traces entirely
	- GOTRACEBACK=sinle (the default) behaves as described above
	- GOTRACEBACK=all adds stack traces for all user-created goroutines
	- GOTRACEBACK=system is like "all" btu adds stack frames for run-time functions and shows goroutines created internally by the run-time.
	- GOTRACEBACK=crash is like "system" but crashes in an OS specific manner instead of exiting.
	For example, on Unix systems, the crash raises SIGABRT to trigger a core dump.
	- GOTRACEBACK=wer is like "crash" but doesn't disable Windows Error Reporting. For historical reasons, the GOTRACEBACK settings 0, 1, 2 are synonyms for none, and system, respectively. The runtime/debug.SetTraceback
		function allows increasing the amount of output at run time, but it cannot reduce the amount below that specified by the environment variable.
7. The GOARCH, GOOS, GOPATH, and GOROOT environment variables complete the set of Go environment variables. They influence the building of Go programs.
	- GOARCH, GOOS, and GOROOT are recorded at compile time and made available by constants or functions in this package but they don't influence the execution of the run-time system.

	SECURITY
1. On Unix platforms, Go's runtime system behaves slightly differently when a binary is setuid/
	setgid or executed with with setuid/setgid-like properties, in order to prevent behaviors. On Linux this is determined by checking for the AT_SECURE flag in the auxiliary vector.

	CONSTANTS
1. const Compiler = "gc"
	Compiler is the name of the compiler toolchain that built the running binary. Known toolchains are:
		gc Also known as cmd/compile
		gccgo The gccgo front end, part of the GCC compiler suite
2. const GOARCH string = goarch.GOARCH
	GOARCH os the running program's architecture target on of 386, amd64, arm, s390x, and so on.
3. const GOOS string = goos.GOOS
	GOOS is the running program's operating system target: one of darwin, freebsd, linux and so on.

	To view possible combinations of GOOS and GOARCH, run "goo tool dist list"

	VARIABLES
1. MemProfileRate int = 512 * 1024
	MemProfileRate controls the fraction of memory allocations that are recorded and reported in the memory profile. The profiler aims to sample an average of one allocation per MemProfileRate
	bytes allocated.

	- To include every allocated block in the profile, set MemProfileRate to 1.
	- To turn off profiling entirely, set MemProfileRate to 0

	The tools that process the memory profiles assume that the profile rate is constant across the lifetime of the program and equal to the current value.

	Programs that change the memory profiling rate should do so just once, as early as possible in the execution of the program (for example, at the beginning of main).

	FUNCTIONS
1. func BlockProfile(p []BlockProfileRecord) (n int, ok bool)
	BlockProfile returns n, the number of records in the current blocking profile.

	If len(p) >= n, BlackProfile copies the profile into p and returns n, true.

	If len(p) < n, BlockProfile doesn't change p and returns n, false.

	Most clients should usee the runtime/pprof package or the testing package's -test.blockprofile flag instead of calling BlockProfile directly.
2. func Breakpoint()
	Breakpoint executes a breakpoint trap
3. func Caller(skip int) (pc uintptr, file string, ok bool)
	Caller reports file and line number information about function invocations on the calling goroutine's stack.

	The argument skip is the number of stack frames to ascend, with 0 identifying the caller of
	Caller. (For historical reasons the meaning of skip differs between Caller and Callers).

	The return values report the program counter, file name, and line number within the file of the
	corresponding call.

	The boolean ok is false if it wasn't possible to recover the information.
4. func Callers(skip int, pc []uintptr) int
	Callers fills the slice pc with return program counters of function invocations on the calling goroutine's stack.

	The argument skip is the number of stack frames to skip before recording in pc, with 0
	identifying the frame for Callers itself and 1 identifying the caller of Callers.

	it returns the number of entries written to pc.

	To translate these PCs into symbolic information such as function names and line numbers, use
	CallerFrames. CallerFrames accounts for inlined functions and adjusts the return program
	counters into call program counters. Iterating over the returned slice of PCs directly is
	discouraged, as is using FuncForPC on any of the returned PCs, since these cannot account for
	inlining or return program counter adjustment.
5. func GC()
	GC runs a garbage collection and blocks the caller until the garbage collector is complete. It
	may also block the entire program.
6. func GOMAXPROCS(n int) n
	GOMAXPROCS sets the maximum number of CPUs that can be executing simultaneously and returns the previous setting.

	It defaults to the value of runtime.NumCPU.

	If n < 1, it doesn't change the current setting.

	This call will go away when the scheduler improves.
7. func GOROOT()
	GOROOT returns the root of the Go tree.

	It uses GOROOT env. variable, if set at process start, or else the root used during the Go
	build.
8. func Goexit()
	Goexit terminates the goroutine that calls it. No other goroutines is affected.

	Goexit runs all deferred calls before terminating the goroutine. Because Goexit is not a panic, any recover calls in those deferred functions will return nil.

	Calling Goexit from the main goroutine terminates that goroutine without func main returning. Since func main hasn't returned, the program continues execution of other goroutines. If all other goroutines exit, the
	program crashes.
9. func GoroutineProfile(p []StackRecord) (n int, ok bool)
	GoroutineProfile returns n, the number of records in the active goroutine stack profile.

	- If len(p) >= n, GoroutineProfile copies the profile into p and returns n, true.
	- If len(p) < 0, GoroutinesProfile doesn't change p and returns n, false

	Most clients should use the runtime/pprof package instead of calling GoroutineProfile directly.
10. func Gosched()
	Gosched yields the processor, allowing other goroutines to run.
	It doesn't suspend the current goroutine, so execution resumes automatically.
11. func KeepAlive(x any)
	KeeapAlive marks its argument as currently reachable. This ensures that the object is not
	freed, and its finalizer isn't run, before the point in the program where KeepAlive is called.

	Note: KeepAlive(x any) should only be used to prevent finalizers from running prematurely, when
	used with unsafe.Pointer, the rules for valid uses of unsafe.Pointer still apply.
12. func LockOSThread()
	LockOSThread wires the calling goroutine to its current operating system thread. The calling
	goroutine will always execute in that thread, and no other goroutine will execute in it, until
	the calling goroutine has made as many calls to UnlockOSThread as to LockOSThread.

	If the calling goroutine exits without unlocking the thread, the thread will be terminated.

	All init functions are run on the startup thread. Calling LockOSThread from an init function
	will cause the main function to be invoked on that thread.

	A goroutine should call LockOSThread before calling OS services or non-Go library functions
	that depends on per-thread state.
13. func MemProfile(p []MemProfileRecord, inuseZero bool) (n int, ok bool)
	MemProfile returns a profile of memory allocated and freed per allocation site.

	MemProfile returns n, the number of records in the current memory profile.
	- If len(p) > n, MemProfile copies the profile into p and returns n, true
	- If len(p) < n, MemProfile doesn't change p and returns n, false

	If inuseZero is true, the profile includes allocation records where r.AllocBytes > 0 but r.
	AllocBytes == r.FreeBytes. These are sites where memory was allocated, but it has all been
	released back to the runtime.

	The returned profile may be up to two garbage collection cycles old. This is to avoid skewing
	the profile toward allocations; because allocations happen in real time but frees are delayed
	until the GC perfomes sweeping, the profie only accounts for allocations that have had a chance to be freed by the GC.

	Most clients should use the runtime/pprof package or the testing package's -test.memprofiling flag instead of calling MemProfile directly.
14. func MutexProfile(p []BlockProfileRecord) (n int, ok bool)
	Mutex profile returns n, the number of records in the current mutex profile.
	- If len(p) >= n, MutexProfile copies the profile into p and returns n, true
	- If len(p) < 0, MutexProfile doesn't change p, and returns n, false

	Most clients should use the runtime/pprof package instead of calling MutexProfile directly.
15. func NumCPU() int
	NumCPU returns the number of logical CPUs usable by the current process.

	The set of available CPUs is checked by querying the OS at process startup.
	Changes to OS CPU allocation after process startup are not reflected.
16. func NumCgoCall() int64
	NumCgoCall() returns the number of cgo calls made by the current process.
17. func NumGoroutine() int
	NumGoroutine returns the number of goroutines that currently exist.
18. func ReadMemStats(m *MemStats)
	ReadMemStats populates m with memory allocator statistics.

	The returned memory allocator statistics are up to date of the call to ReadMemStats. This is in
	contrast with a heap profile, which is a snapshot as of the most recently completed garbage
	collection cycle.
19. func ReadTrace() []byte
	ReadTrace returns the next chunk of binary tracing data, blocking until data is available. If
	tracing is turned off and all the data is accumulated while it was on has been returned,
	ReadTrace returns nil.

	The caller must copy the returned data before calling ReadTrace again.

	ReadTrace must be called from one goroutine at a time.

...

20. func Stack(buf []byte, all bool) int
	Stack formats a stack trace of the calling goroutine into buf and returns the number of bytes
	written to buf. If all is true, Stack formats stack traces of all other goroutines into buf
	after the trace for the current goroutine.
21. func StartTrace() error
	StartTrace enables tracing for the current process. While tracing, the data will be buffered
	and available via ReadTrace. StartTrace returns an error if tracing is already enabled. Most
	clients should use the runtime/trace package or testing package's -test.trace flag instead of
	calling StartTrace directly.
22. func StopTrace()
	StopTrace stops tracing, if it was previously enabled. StopTrace only returns after all the
	reads for the trace have completed.
23. func ThreadCreateProfile(p []StackRecord) (n int, ok bool)
	ThreadCreateProfile returns n, the number of records in the thread creation profile.
	- If len(p) >= n, ThreadCreateProfile copies the profile into p and returns n, true
	- If len(p) < n, ThreadCreateProfile doesn't change p and returns n, false

	Most clients should use the runtime/pprof package instead of calling ThreadCreateProfile
	directly
24. func UnlockOSThread()
	UnlockOSThread() undoes an earlier call to LockOSThread. If this drops the number of active
	LockOSThread calls on the calling goroutine to zero, it unwires the calling goroutine from its
	fixed OS Thread. If there are no active LockOSThread calls, this is a no-op.

	Before calling UnlockOSThread, the caller must ensure that the OS Thread is suitable for
	running other goroutines. If the caller made any permanent changes to the state of the thread
	that would affect other goroutines, it shouldn't call this function and thus leave the
	goroutine locked to the OS thread until the goroutine (and hence the thread) exits.
25. func Version() string
	Version returns the Go tree's version string. It's either the commit hash and date at the time
	of the build of, when possible, a release tag like "go1.3".
*/

func main() {

}
