package e_functions

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/html"
)

func SimpleFailureFlow() (string, bool) {
	namesCount := map[string]int{"Dima": 1, "Denis": 3}

	nameCount, ok := namesCount["Diana"]
	if !ok {
		return fmt.Sprintf("Any Dianas are out of the name list"), false

	}
	return fmt.Sprintf("Here's %d Diana(s)\n", nameCount), true
}

func PropagatingSubroutineFailure(url string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return nil
}

func DescriptiveAdditionalInformation(url string) (*html.Node, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("while the requesting to %s: %s", url, err)
	}
	defer response.Body.Close()

	doc, err := html.Parse(response.Body)
	if err != nil {
		return nil, fmt.Errorf("while parsing %s as HTML: %s", url, err)
	}

	return doc, nil
}

func ServerResponseHandler(url string, serverResultError error) {
	if err := WaitForServer(url); err != nil {
		// fmt.Fprintf(os.Stderr, "Site is down: %v\n", err)
		log.Fatalf("Site is down: %v\n", err) // The alternative way to stop the program.
		os.Exit(1)
	}
}

func WaitForServer(url string) error {
	// Set a timeout
	const (
		timeout = 1 * time.Minute
	)
	var deadline = time.Now().Add(timeout)

	for tries := 0; time.Now().Before(deadline); tries++ {
		_, err := http.Head(url)
		if err == nil {
			return nil // success
		}
		log.Printf("server isn't responding: (%s); retrying...", err)
		// This state doubles latency up after each unsuccessful request.
		time.Sleep(time.Second << uint(tries)) // exponential back-off
	}
	return fmt.Errorf("server %s failed to repond after %s", url, timeout)
}

func DiscardingError() error {
	dir, err := os.MkdirTemp("", "scratch")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %v", err)
	}

	fmt.Println(dir)

	os.RemoveAll(dir) // If the directory doesn't exist, the function does nothing and just returns nil.
	// On the other hand if any errors occur, the function ignores them.
	return nil
}

func HandleEOF() {
	input := bufio.NewReader(os.Stdin)
	for {
		r, _, err := input.ReadRune()
		if err == io.EOF {
			fmt.Println(err)
			break
		}
		fmt.Println(r)
	}
}
