package chapter8

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
	"unicode"
)

const (
	serverPrefix = "localhost:"
	trimSet      = "\t\n\v\f\r \u0085\u00A0"
)

var (
	port                 = flag.String("port", "", "port to read time of the particular timezone")
	validPromptArgsCount int
)

func setAddress(server interface{}) error {

	validatedPort, err := isPortValid()
	if err != nil {
		return fmt.Errorf("validating port; %s", err)
	}

	address := serverPrefix + validatedPort

	switch server := server.(type) {
	case *TimeServer:
		server.address = address
	case *ServerFTP:
		server.address = address
	default:
		return fmt.Errorf("unexpeted server type; type: %T", server)
	}

	return nil
}

func parsePortNumber(str string) (string, error) {

	// Check the characters validity
	for _, character := range str {
		if !unicode.IsDigit(character) {
			return "", fmt.Errorf("port consist other characters than digits; typed port: \"%s\"", str)
		}
	}
	return str, nil
}

func isPortBusy(port string) error {
	// Check whether the port is already used
	testConnection, err := net.Dial("tcp", serverPrefix+port)
	// If we can't establish connection, it means that the port is free
	if err != nil {
		return nil
	}
	defer testConnection.Close()

	// Otherwise the connection, the port is already used.
	return fmt.Errorf("port is already used; port: %s", port)
}

func setTimeLocation(server interface{}) error {
	if server == nil {
		return fmt.Errorf("given time server instance is nil")
	}

	// Set the default value of the timezone
	timezone := "UTC"

	// Extract the environment variable's value
	if tz := os.Getenv("TZ"); tz != "" {
		timezone = tz
	}

	// Set the location fiel
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return err
	}

	switch server := server.(type) {
	case *TimeServer:
		server.location = loc
	case *ServerFTP:
		server.location = loc
	default:
		return fmt.Errorf("unexpeted server type; type: %T", server)
	}

	return nil
}

func setListener(server interface{}) error {
	const networkType = "tcp"
	var err error

	switch server := server.(type) {
	case *TimeServer:
		server.listener, err = net.Listen(networkType, server.address)
	case *ServerFTP:
		server.listener, err = net.Listen(networkType, server.address)
	default:
		return fmt.Errorf("unexpeted server type; type: %T", server)
	}

	if err != nil {
		return err
	}

	return nil
}

func isPortValid() (string, error) {
	flag.Parse()

	validPromptArgsCount = 2

	// Check whether the args count is valid
	if len(os.Args) < validPromptArgsCount {
		return "8000", nil
	}

	// Check whether the port consists only of digits
	validatedPort, err := parsePortNumber(*port)
	if err != nil {
		return "", fmt.Errorf("port value consists not only of digits; typed value: %s", err)
	}

	return validatedPort, nil
}

// userInput() Forces the user to write in the correct ftp methods and its arguments for ftp work
func userInput() ([]string, error) {
	const maxAttemptsCount = 5
	var (
		scanner = bufio.NewScanner(os.Stdin)

		userInput    []string
		attemptCount int
	)
	scanner.Split(bufio.ScanLines)

	// Get values for parameters from user
	for scanner.Scan() {
		fmt.Print("Write ftp method and arguments for it separated by commas -> ")
		if attemptCount == maxAttemptsCount {
			return nil, fmt.Errorf("attempts limit exceeded")
		}

		userInput = strings.Split(strings.Trim(scanner.Text(), trimSet), ",")

		if err := isInputValid(userInput); err != nil {
			fmt.Printf("%s\n", err)
			attemptCount++
			continue
		}

		break
	}

	return userInput, nil
}

func formatParams(params []string) string {
	result := "| "

	for _, param := range params {
		result += param + " "
	}

	result += "|"

	return result
}

func isInputValid(args []string) error {
	const errorPrefix = "invalid input: "

	if len(args) < 1 {
		return fmt.Errorf(errorPrefix + "no method written")
	}

	for _, arg := range args {
		if len(arg) == 0 {
			return fmt.Errorf(errorPrefix + "empty argument")
		}
	}

	return nil
}
