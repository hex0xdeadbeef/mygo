package contextpkg

import (
	"context"
	"fmt"
	"sync"
	"time"
)

/*
With the "context" package we can communicate extra information alongside the simple notification to cancel: why the cancellation was occuring, or whether or not our
function has a deadline by which it needs to complete.
*/

/*
If we use the "context" package, each function that is downstream from our top-level concurrent call would take in a "Context" as its first argument.
*/

/*
The purpose of the "Value()" function is passing request-specific information along in addition to information about preemption.
*/

/*
The "context" package serves two primary purposes:
	1) To provide an API for cancelling branches of your call graph
	2) To provide a data-bag for transporting request-scoped data through your call graph
*/

/*
The "context" package helps to serve the following ascpects:
	1) A goroutine's parent may want to cancel it
	2) A goroutine may want to cancel its children
	3) Any blocking operations within a goroutine need to be preemptable so that it may be canceled.
*/

/*
We mustn't save the references of "context.Context" instances. Instances of "context.Context" may look equivalent from the outside, but internally they may change at
every stack-frame. For this reason it's important to always pass instances of Context into your functions. This way functions have the Context inteded for it, and
not the Context intended for a stack-frame N levels up the stack.
*/

/*
At the top of our asynchronous call-graph, our code probably won't have been passed a Context. To start the chain, the "context" package provides us with two functions
to create empty instances of Context.
	func Background() Context
	func TODO() Context - this function's purpose is to serve as a placeholder for when we don't know which Context to utilize, or we expect our code to be provided
	with a Context, but the upstream code hasn't yet furnished one.
*/

/* func DoneUsing() {
	var (
		wg   sync.WaitGroup
		done = make(chan struct{})
	)
	defer close(done)

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := printGreetings(done); err != nil {
			fmt.Printf("%v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := printFarewell(done); err != nil {
			fmt.Printf("%v", err)
			return
		}

	}()

	wg.Wait()
}

func printGreetings(done <-chan struct{}) error {
	greeting, err := genGreeting(done)
	if err != nil {
		return err
	}

	fmt.Printf("%s world!\n", greeting)

	return nil
}

func printFarewell(done <-chan struct{}) error {
	farewell, err := genFarewell(done)
	if err != nil {
		return err
	}

	fmt.Printf("%s world!\n", farewell)

	return nil
}

func genGreeting(done <-chan struct{}) (string, error) {
	switch locale, err := locale(done); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "hello", nil
	}

	return "", fmt.Errorf("unsupported locale")
}

func genFarewell(done <-chan struct{}) (string, error) {
	switch locale, err := locale(done); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "goodbye", nil
	}

	return "", fmt.Errorf("unsupported locale")
}

func locale(done <-chan struct{}) (string, error) {
	select {
	case <-done:
		return "", fmt.Errorf("canceled")
	case <-time.After(3 * time.Second):
	}

	return "EN/US", nil

}

*/

func ContextUsing() {
	var (
		wg sync.WaitGroup
		// Create the parent context instance and the prior cancel function
		ctx, cancel = context.WithCancel(context.Background())
	)
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := printGreetings(ctx); err != nil {
			fmt.Printf("%v\n", err)
			// Cancel the parent context if the children have been canceled
			cancel()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := printFarewell(ctx); err != nil {
			fmt.Printf("%v\n", err)
			return
		}

	}()

	wg.Wait()
}

func printGreetings(ctx context.Context) error {
	greeting, err := genGreeting(ctx)
	if err != nil {
		return err
	}

	fmt.Printf("%s world!\n", greeting)

	return nil
}

func printFarewell(ctx context.Context) error {
	farewell, err := genFarewell(ctx)
	if err != nil {
		return err
	}

	fmt.Printf("%s world!\n", farewell)

	return nil
}

func genGreeting(ctx context.Context) (string, error) {
	// Wrap the parent context to be passed downward into the childs and cancel any children after 1 second.
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	switch locale, err := locale(ctx); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "hello", nil
	}

	return "", fmt.Errorf("unsupported locale")
}

func genFarewell(ctx context.Context) (string, error) {
	switch locale, err := locale(ctx); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "goodbye", nil
	}

	return "", fmt.Errorf("unsupported locale")
}

func locale(ctx context.Context) (string, error) {
	if deadline, ok := ctx.Deadline(); ok {
		if deadline.Sub(time.Now().Add(1*time.Minute)) <= 0 {
			return "", context.DeadlineExceeded
		}
	}

	select {
	case <-ctx.Done():
		// Return the reason of cancellation
		return "", ctx.Err()
	case <-time.After(3 * time.Second):
	}

	return "EN/US", nil

}

/*
Context with the value
*/

/*
Qualifications about context.WithValue:
	1) The key we use must satisfy Go's notion of comparability. That is, the equality operators == and != need to return correct results when used.
	2) Values returned must be safe to access from multiple goroutines.
	3) We must wrap the underlying comparable types of keys into our own types.
*/

/*
The heuristics about what we can store in the context.Context:
	1) The data should transit process or API boundaries
	2) The data should be immutable
	3) The data should trend toward simple types
	4) The data should be data, not types with methods
	5) The data should help decorate operations, not driven them.
	6) How many layers this data might need to traverse before utilization.
*/

func ContextWithDataUsing() {

	ProcessRequest("hex0xdead", "akld319qdjo17")
}

type ctxKey int

const (
	ctxUserID ctxKey = iota
	ctxAuthToken
)

func UserID(ctx context.Context) string {
	return ctx.Value(ctxUserID).(string)
}

func AuthToken(ctx context.Context) string {
	return ctx.Value(ctxAuthToken).(string)
}

func ProcessRequest(userID, authToken string) {
	ctx := context.WithValue(context.Background(), ctxUserID, userID)
	ctx = context.WithValue(ctx, ctxAuthToken, authToken)

	HandleResponse(ctx)
}

func HandleResponse(ctx context.Context) {
	fmt.Printf("handling response for %v (%v)", UserID(ctx), AuthToken(ctx))
}
