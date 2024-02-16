package chapter8

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	readingDurationMillis  = 100
	serverReaderBufferSize = 8192
)

var serverLog *log.Logger

func init() {
	var err error
	serverLog, err = NewLogger(serverLogFile)
	if err != nil {
		log.Fatalf("creating server logger; %s", err)
	}
}

type ServerFTP struct {
	address  string
	location *time.Location
	listener net.Listener

	usersCredentials map[string]string
	mu               sync.Mutex
}

func NewServerFTP() (*ServerFTP, error) {
	newServerFTP := ServerFTP{
		usersCredentials: map[string]string{"hex0xdead": "abc", "hex0xbeef": "efg"},
	}

	if err := setTimeLocation(&newServerFTP); err != nil {
		return nil, fmt.Errorf("setting the timezone; %s", err)
	}

	if err := setAddress(&newServerFTP); err != nil {
		return nil, fmt.Errorf("setting server address; %s", err)
	}

	if err := setListener(&newServerFTP); err != nil {
		return nil, fmt.Errorf("setting listener; %s", err)
	}

	return &newServerFTP, nil
}

func (ts *ServerFTP) Start() {
	for {
		// Accept() blocks the programm execution until an incoming connection request is made, then returns a net.Conn object
		// representig the connection
		connection, err := ts.listener.Accept()
		if err != nil {
			log.Print(err) // e.g. conncetion aborted
		}

		go ts.handleClient(connection) // handle connections concurrently
	}
}

func (ts *ServerFTP) handleClient(conn net.Conn) {
	defer conn.Close()
	connReaderWriter := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	// Create the buffer to read data from the connection
	var buf []byte = make([]byte, serverReaderBufferSize)

	// Try to read client's request data
	n, err := connReaderWriter.Read(buf)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			log.Printf("reading connection data: %s", err)
		} else {
			log.Printf("reading connection data: %s", err)
		}
	}

	// Reduce the buffer to the client's size of the data
	request := strings.Fields(string(buf[:n]))
	response := ""

	switch request[0] {
	case "ls":
		// Get entries in the current directory
		entries, err := os.ReadDir(".")
		if err != nil {
			response += fmt.Sprintf("internal error", err)
			break
		}

		namesInLine := 0
		// Put all the filenames into the string
		for _, entry := range entries {
			info, err := entry.Info()
			if err != nil {
				response += fmt.Sprintf("internal error", err)
				break
			}
			if namesInLine == 3 {
				namesInLine = 0
				response += "\n"
			}
			response += info.Name() + " "
			namesInLine++
		}

	case "cd":
		// request[1] - path to the new directory
		err := os.Chdir(request[1])
		if err != nil {
			response += fmt.Sprintf("internal error", err)
			break
		}
		workDirectory, err := os.Getwd()
		if err != nil {
			response += fmt.Sprintf("internal error", err)
			break
		}
		response += workDirectory
	case "get":
		newFile, err := os.OpenFile(request[1], os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			response += fmt.Sprintf("internal error", err)
			break
		}
		defer newFile.Close()

		fileWriter := bufio.NewWriter(newFile)
		fileWriter.WriteString(strings.Join(request[2:], " "))
		if err := fileWriter.Flush(); err != nil {
			log.Printf("flushing the writer's buffer: %s", err)
		}

	default:
		response += "400 Bad Request"
	}

	connReaderWriter.Write([]byte(response))
	if err := connReaderWriter.Flush(); err != nil {
		log.Printf("flushing the reader's buffer: %s", err)
	}

}
