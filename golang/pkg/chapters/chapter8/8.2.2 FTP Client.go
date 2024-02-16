package chapter8

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

const (
	clientReaderBufSize = 1024
)

var (
	validMethods map[string]int = map[string]int{"cd": 1, "ls": 0, "get": 2, "close": 0}
)

type Client struct {
	conn net.Conn
}

// NewClient() creates a new client if all the checks is passed, returning the clien and non-nil err, otherwise returns nil and the
// corresponding non-nil error.
func NewClient() (*Client, error) {
	flag.Parse()
	validPromptArgsCount = 2

	// Check whether the args count is valid
	if len(os.Args) != validPromptArgsCount {
		return nil, fmt.Errorf("invalid args count; args count: %d", len(os.Args))
	}

	// Check the validity of the port
	validatedPort, err := isPortValid()
	if err != nil {
		return nil, fmt.Errorf("validating port; %s", err)
	}

	// Try to establish connection
	address := serverPrefix + validatedPort
	connection, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("establishing connection; %s", err)
	}

	// Create a client with a connection
	client := Client{conn: connection}

	return &client, nil
}

func (c *Client) StartSession() {
	go c.sendCommands()

	for {
		c.getAnswer()
	}
}

// (c *Client) sendCommands() Writes the message into the connection
func (c *Client) sendCommands() {
	var connWriter bufio.Writer = *bufio.NewWriter(c.conn)

	for {
		commAndArgs, err := c.constructCommand()
		if err != nil {
			log.Printf("constructing command: %s", err)
			continue
		}

		switch commAndArgs[0] {
		case "cd":
			connWriter.WriteString(strings.Join(commAndArgs, " "))
		case "ls":
			connWriter.WriteString(strings.Join(commAndArgs, " "))
		case "get":

		case "close":
			connWriter.WriteString(strings.Join(commAndArgs, " "))
		}

		if err := connWriter.Flush(); err != nil {
			log.Printf("flushing buffer: %s", err)
		}

	}
}

// (c *Client) constructCommand() Gets an user input, validates the input, sends to the caller the splitted input
func (c *Client) constructCommand() ([]string, error) {
	commAndArgs, err := userInput()
	if err != nil {
		return nil, fmt.Errorf("constructing command: %s", err)
	}

	if err := c.validateCommAndArgs(commAndArgs); err != nil {
		return nil, err
	}

	return commAndArgs, nil
}

// (c *Client) validateCommAndArgs Checks the validity of the typed by user data
func (c *Client) validateCommAndArgs(commAndArgs []string) error {
	_, ok := validMethods[commAndArgs[0]]

	if !ok {
		return fmt.Errorf("invalid ftp method: written method: \"%s\"", commAndArgs[0])
	}

	if len(commAndArgs[1:]) != validMethods[commAndArgs[0]] {
		return fmt.Errorf("invalid args count for \"%s\" method: args count: %d", commAndArgs[0], len(commAndArgs[1:]))
	}

	return nil
}

// (c *Client) getAnswer() Checks the server's response and prints it out into the standart output
func (c *Client) getAnswer() error {
	// Set limit on reading the server response time
	if err := c.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 500)); err != nil {
		return fmt.Errorf("setting reading deadline: %s", err)
	}

	// Create the reader for servers responses
	connReader := bufio.NewReader(c.conn)
	var buf []byte = make([]byte, clientReaderBufSize)

	// Try to read server response
	n, err := connReader.Read(buf)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return fmt.Errorf("reading connection data: %s", err)
		} else {
			return fmt.Errorf("reading connection data: %s", err)
		}
	}

	// If everything is okay, reduce the buffer and write out a response into the standart output file
	buf = buf[:n]
	fmt.Printf("Server answer: %s", string(buf))
	return nil
}
