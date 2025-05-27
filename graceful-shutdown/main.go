package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

/*
	GRACEFUL SHUTDOWN

	Graceful shutdown in any application generally satisfies three minimum conditions:
	1. Close the entry point by stopping new requests or messages from source like HTTP, pub/sub
	systems, etc. However, keep outgoing conns to third-party services like databases or caches alive.
	2. Wait for all ongoing requests to finish. If a request takes too long, respond with a graceful error.
	3. Release critical resources such as database conns, file locks, or network listeners.

*/

/*
	1. CATCHING THE SIGNAL

	Before we handle graceful shutdown, we first need to catch termination signals. These signals
	tell our application that `it's time to exit` and begin the shutdown process.

	So, what are signals?

	In Unix-like systems, signals are software interrupts. They notify a process that something has
	happened and it should take action. When a singal is sent, the OS interrputs the normal flow of the process to deliver the notification.

	Here are a few possible behaviors:
		- Signal handler: A process register a handler (a function) for a specific signal. This
		function runs when that signal is received.
		- Default action: If no handler is registered, the process follows the default behavior for
		that signal. This might mean terminating, stopping, continuing, or ignoring the process.
		- Unblockable signals: Some signals, like `SIGKILL` (signal number 9), cannot be caught or
		ignored. They may terminate the process.

	When our Go application starts, even before our main function runs, the Go runtime automatically registers signal handlers for many signals (SIGTERM, SIGQUIT, SIGILL, SIGTRAP and others). However, for graceful shutdown, only three termination signals are typically important:
		- SIGTERM (Termination): A standard and polite way to ask the process to terminate. It
		doesn't force the process to stop. Kubernetes sends this signal when it wants our application to exit before it forcibly kills it.
		- SIGINT (Interrupt): Sent when the user wants to stop a process from the terminal, usually by pressing Ctrl + C
		- SIGHUP (Hang up): Originally used when a terminal disconnected. Now, it's often repurposed to signal an application to reload its configuration.
	People mostly care about SIGTERM and SIGINT, SIGHUP is less used today for shutdown and more for reloading configs.

	By default, when out application receives SIGTERM and SIGINT, or SIGHUP the GO runtime will terminate the application.

	NOTE: When our Go app gets a SIGTERM, the runtime first catches it using a built-in handler. It checks if a custom handler is registereed. If not, the runtime disables its own handler temporarily, and sends the same signal (SIGTERM) to the application again. This time. the OS handles it using the default behavior, which is to terminate the process.

	We can override it by registering our own signal handler using the `os/signal` package.

	signal.Notify tells the Go runtime to deliver specified signals to a channel instead of using the default behavior. This allows us to handle them manually and prevents the app from terminating automatically.

	A buffered channel with a capacity 1 is a good choice for reliable signal handling. Internally Go sends signals to this channel using a `select` statement with a default case:
		select {
			case c <- sig:
			default:
		}

	This is different from the usual `select` used with receiving channels. When used for sending:
		- If the buffer has space, the signal is sent and the code continues
		- If the buffer is full, the signal is discarded, and the `default` case runs. If we're using an unbuffered channel and no goroutines is actively receiving, the signal will be missed.

	Even though it can only hold one signal, this buffered channel helps avoid missing that first signal while our app is still initializing and not yet listening.

	NOTE: We can call signal.Notify multiple times for the same signal. Go will send that signal to all registered channels.

	When we press Ctrl + C more than once, it doesn't automatically kill the app. The first Ctrl + C sends a SIGINT to the foreground process. Pressing it again usually sends another SIGINT, not SIGKILL. Most terminals, like bash or other Linux shells, don't escalate the signal automatically. if we want to force a stop, we must send SIGKILL manually using `kill -9`

	This isn't ideal for local development, where we may want the second Ctrl + C to terminate the app forcefully. We can stop the app from listening to further signals by using signal.Stop right after the first signal is received.

	Starting with Go 1.16, we can simplify signal handling by using `signal.NotifyContext`, which ties signal handling to context cancelation.
*/

/*
	2. TIMEOUT AWARENESS
	It's important to know how long our app has to shut down after receiving a termination signal. For example, in k8s, the default grace period is 30 seconds, unless otherwise specified using terminationGracePeriodSeconds field. After this period, k8s sends a SIGKILL to forcefully stop the app. This signal cannot be caught or handled.

	Our shutdown logic must complete within this time, including processing any remaining reqs and releasing resources.

	Assume, the default is 30 secs. It's a good practice to reserve about 20 percent of the time as a safety margin to avoid being killed before cleanup finishes. This means aiming to finish everything within 25 seconds to avoid data loss or inconsistency.
*/

/*
	3. STOP ACCEPTING NEW REQS
	When using `net/http` we can handle graceful shutdown by calling the http.Server.Shutdown method. This method stops the server from accepting new conns and waits for all active reqs to complete before shutting down idle conns.

	Here's how it behaves:
		- If a req is already in progress on an existing conn, the server will allow it to complete. After that, the conn is marked as idle and closed.
		- If a client tries to make a new conn during shutdown, it'll fail because the server's listeners are already closed. This typically results in a "conn refused" error.

	In a containerized environment (and many others orchestrated environments with external load balancers), don't stop accepting new requests immediately. Even after a pod is marked for termination, it might still receive traffic for a few moments.

	k8s internal components like kube-proxy are quickly aware of the change in pod status to `Terminating`. They then prioritize routing internal traffic to Ready, Serving endpoints over `Terminating, Serving` ones.

	The external load balancer, however, operates independently from k8s. It typically uses its own health check mechanisms to determine which backend nodes should receive traffic. This healthcheck indicates whether there are healthy (Ready) and non-terminating pods on the node. However, this check needs a little time to propagate. And there are two ways to handle this:
		- Use a preStop hook to sleep for a while, so the external load balancer has time to recognize
		that the pod is terminating:
			lifecycle:
			prestop:
				exec:
				command: ["/bin/sh", "-c", "sleep 10"]

		And really importantly, the time taken by the preStop hook is included within the terminationGracePeriodSeconds.

		- Fail the readiness probe and sleep at the code level. This approach isn't only applicable
		to k8s environments, but also to other environments with load balancers that need to know
		the pod is not ready.

		What is readiness probe?
			A readiness probe determines when a container is prepared to accept traffic by periodically checking its health through configured methods like HTTP reqs, TCP conns, or command executions. If the probe fails, k8s removes the pod from the service's endpoints, preventing it from receiving traffic until it becomes ready again.

		To avoid conn errors during this short window, the correct strategy is to fail the readiness probe first. This tells the orchestrator that our pod shouldn't receive traffic.

		This pattern is also used as a code example in the test images. In their implementation, a closed channel is used to signal the readiness probe to return HTTP 503 when the app is preparing to shut down.

		After updating the readiness probe to indicate that the pod is no longer ready, wait a few seconds to allow the system time to propagate the change.

		The exact wait time depends on our readiness probe config; we'll use 5 seconds for this article with the following simple config.
			readinessProbe:
				httpGet:
					path: /healthz
					port: 8080
				periodSeconds: 5
*/

/*
	4. HANDLE PENDING REQS

	Now that we're using shutting down the server gracefully, we need to choose a timeout based on our shutdown budget:
		ctx, cancelFn := context.WithTimeout(context.TODO, timeout)
		err := server.Shutdown(ctx)
	The server.Shutdown func returns in only two situations:
		- All active conns are closed and all handlers have finished processing
		- The context passed to Shutdown(ctx) expires before the handlers finish. In this case, the server gives up waiting.
	In either case, Shutdown only returns after the server has completely stopped handling reqs. This is why handlers must be fast and context-aware. Otherwise, they may be cut off mid-process in case 2, which can cause isssues like partial writes, data loss, inconsistent, state, open txs, or corrupted data.

	A common issue is that handlers aren't automatically aware when the server is shutting down.

	So, how can we notify our handlers that the server is shutting down? The answer is by using context. There are two main ways to do this:
		- Use context middleware to inject cancellation logic. This middleware wraps each req with a context that listens to a shutdown signal.
		- Use BaseContext to provide a global context to all connections. In an HTTP server, we can customize two types of contexts: BaseContext and ConnContext. For GS, BaseContext is more suitable. It allows us to create a global context with cancellation that applies to the entire server, and we can cancel it to signal all active reqs that the server is shutting down.

	All of this work around GS won't help if our funcs don't respect context cancellation. Try to avoid using context.Background(), time.Sleep(), or any other function that ignores context.

	For example, time.Sleep(duration) can be replaced with a context-aware version like this one:
		func Sleep(ctx context.Context, dur time.Duration) error {
			select {
			case <-time.After(dur):
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		}

	The core principle of GS is the same across all systems:
		Stop accepting new reqs or messages, and give existing ops time to finish within a defined grace period.

	Some may wonder about the server.Close() method, which shuts down the ongoing conns immediately without waiting for reqs to finish. Can it be used after server.Shutdown() returns an error?

	The short answer is yes, but it depends on your GS strategy. The Close() method forcefully closes all active listeners and conns:
		- Handlers that are actively using the network will receive errors when they try to read or write.
		- The client will immediately receive a conn error, such as ECONNRESET (`socket hang up`)
		- However, long-running handlers that aren't interacting with the network may `continue running` in the background.

	This is why using context to propagate a shutdown signal is still the more reliable and graceful approach.
*/

/*
	5. RELEASE CRITICAL RESOURCES
	A common mistake is releasing critical resources as soon as the termination signal is received. At that point, your handlers and in-flight reqs may still be using those resources. You should delay the resource cleanup until GS timeout has passed or all reqs are done.

	In many cases, simply letting the process exit is enough, The OS will automatically reclaim resources. For instance:
		- Memory allocated by Go is automatically freed when the process terminates
		- File descriptors are closed by the OS
		- OS-level resources like process handles are reclaimed

	However, there are important cases where explicit cleanup is still necessary during GS:
		- Database conns should be closed properly. If any txs are still open, they need to be commited or rolled back. Without a proper shutdown, the DB has to rely on conn timeouts.
		- Message queues and brokers often require a clean shutdown. This may involve flushing messages, commiting offsets, or signaling to the broker that the client is exiting. Without this, there can be rebalancing issues or message loss.
		- External services may not detect the disconnect immediately. Closing conns manually allows those systems to clean up faster than waiting for TCP timeouts.

	A good rule is to shut down components in the reverse order of how they were initialized. This respects dependencies between components.

	We should use defer statements to do it.
*/

func ExampleA() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Setup work here

	<-sigChan

	fmt.Println("Received termination signal, shutting down...")
}

func ExampleB() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT)

	<-sigCh

	// Setup work here

	signal.Stop(sigCh)
	select {}
}

func ExampleC() {
	ctx, stop := signal.NotifyContext(context.TODO(), syscall.SIGINT, syscall.SIGTERM)
	// We should still call stop after ctx.Done() to allow a second Ctrl + C to forcefully terminate the app.
	defer stop()

	// Setup tasks here

	<-ctx.Done()
	stop()
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	if isShuttingDown.Load() {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("shutting down"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func WithGracefulShutdown(next http.Handler, cacnelCh <-chan struct{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := WithCancellation(r.Context(), cacnelCh)
		defer cancel()

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})

}

func WithCancellation(ctx context.Context, ch <-chan struct{}) (context.Context, context.CancelFunc) {
	return context.TODO(), func() {}
}

func BaseContextPropagation() {
	ongoingCtx, cancelFn := context.WithCancel(context.Background())

	server := &http.Server{
		Addr:    ":8080",
		Handler: nil,
		BaseContext: func(l net.Listener) context.Context {
			return ongoingCtx
		},
	}

	_ = server

	// After attempting graceful shutdown:
	cancelFn()
	time.Sleep(5 * time.Second) // optional delay to allow context propagation
}

func Sleep(ctx context.Context, dur time.Duration) error {
	select {
	case <-time.After(dur):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
