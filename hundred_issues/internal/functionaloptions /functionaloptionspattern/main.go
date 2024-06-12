package main

import (
	"fmt"
	"net/http"
)

/*
	FUNCTIONAL OPTIONS PATTERN REALIZATION
1. The idea is simple: unexported structure includes the configuration: "options"
	1) Each of its options represent a function itself returning an error of the same type. For example: WithPort
		type Option func(options *options) error
2. Every field of a configuration demands its own public function (It should start with "With" prefix) that includes the analogous logic: a check of input data, changing a
state of internal structure of the configuration.
3. Building a target object we apply all the closures and finally constructs an object based on the resulting options instance
4. gRPC pkg uses the pattern of FUNCTIONAL OPTIONS
*/

func main() {
	NewServer("localhost")

}

// The schema of a function that changes a state of an options instance
type Option func(o *options) error

// The structure of config. It's used to build a targe object based on the values of options.
type options struct {
	port *int
}

// WithPort returns a closure changing an internal state of an options ob. WithPort is the example of a building func.
func WithPort(port int) Option {

	return func(o *options) error {
		if port < 0 {
			return fmt.Errorf("given port %d is negative", port)
		}

		o.port = &port
		return nil
	}

}

// NewServer instantiates a server with closures in the variadic param opts sequentially apllying them over options variable. Finally it builds a server based on options
// applied.
func NewServer(addr string, opts ...Option) (*http.Server, error) {
	var (
		srvr    *http.Server
		options options
	)

	for _, opt := range opts {
		if err := opt(&options); err != nil {
			return nil, fmt.Errorf("applying: %#v; %w", opt, err)
		}
	}
	// some logic ..

	return srvr, nil
}

func NewServerUsageParameterized() {
	srvr, err := NewServer("localhost", WithPort(8090))
	if err != nil {
		panic(err)
	}

	if err := srvr.Close(); err != nil {
		panic(err)
	}
}

// NewServerUsageEmpty instantiates a server based on the zeroed options val
func NewServerUsageEmpty() {
	srvr, err := NewServer("localhost")
	if err != nil {
		panic(err)
	}

	if err := srvr.Close(); err != nil {
		panic(err)
	}
}
