package main

import (
	"math/rand"
	"net/http"
)

/*
	BUILDER PATTERN REALIZATION
1. Constructing a Config object is separated from the structure itself. We use the additional struct ConfigBuilder that is given methods to configure and create a Config file
2. The ConfigBuilder includes a client's configuration. It provides us the method Port to manage port. Usually this type of methods returns the Builder itself, so that we can
use a chain of methods. For example: builder.Foo("foo").Bar("bar").
3. This way complicates error handling
*/

const (
	defaultHTTPPort = 8080
)

type Config struct {
	Port int
}

// ConfigBuilder includes methods to create a Config instance and the pointer field to make decisions while creation
type ConfigBuilder struct {
	port *int
}

// ConfigBuilder returns a pointer to a zeroed ConfigBuilder instance
func EmptyConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{}
}

// Port returns the ConfigBuilder instance with the port field set to the given port argument
func (cb *ConfigBuilder) Port(port int) *ConfigBuilder {
	if cb == nil {
		return &ConfigBuilder{port: &port}
	}

	cb.port = &port
	return cb
}

// Build is the final method in final instance creation. It's used to instantiate the object finally after the chain of creating methods.
func (cb *ConfigBuilder) Build() (cfg Config, err error) {
	cfg = Config{}

	switch port := cb.port; {
	case port == nil:
		cfg.Port = defaultHTTPPort
		break
	case *port == 0:
		cfg.Port = randomPort()
		return
	case *port < 0:
		return cfg, err
	default:
		cfg.Port = *cb.port
	}

	return cfg, nil
}

// randomPort returns a random ports. It's used when the initialized port is set to zero val
func randomPort() int {
	const (
		pad = 1000
	)

	return rand.Intn(9000) + pad
}

// NewServer creates a new server based on args given and returns it and an error if any
func NewServer(addr string, config Config) (*http.Server, error) {
	var (
		srvr *http.Server
	)

	// some logic

	return srvr, nil
}

// NewServerUsage uses the ConfigBuilder in context of http server creation
func NewServerUsage() {
	builder := EmptyConfigBuilder().Port(8090)

	cfg, err := builder.Build()
	if err != nil {
		panic(err)
	}

	server, err := NewServer("localhost", cfg)
	if err != nil {
		panic(err)
	}

	if err := server.Close(); err != nil {
		panic(err)
	}
}
