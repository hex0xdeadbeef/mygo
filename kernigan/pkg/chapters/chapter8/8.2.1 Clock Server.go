package chapter8

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

/*
cd /Users/dmitriymamykin/Desktop/goprojects/golang/cmd

TZ=US/Eastern go run main.go -port 8010 &
TZ=Asia/Tokyo go run main.go -port 8020 &
TZ=Europe/London go run main.go -port 8030 &
TZ=America/Los_Angeles go run main.go -port 8040 &
TZ=Australia/Sydney go run main.go -port 8050 &
*/

type TimeServer struct {
	address  string
	location *time.Location

	listener net.Listener
}

func NewTimeServer() (*TimeServer, error) {

	newTimeServer := TimeServer{}
	if err := setTimeLocation(&newTimeServer); err != nil {
		return nil, fmt.Errorf("setting the timezone; %s", err)
	}

	if err := setAddress(&newTimeServer); err != nil {
		return nil, fmt.Errorf("setting server address; %s", err)
	}

	if err := setListener(&newTimeServer); err != nil {
		return nil, fmt.Errorf("setting listener; %s", err)
	}

	return &newTimeServer, nil
}

func (ts *TimeServer) Start() {
	for {
		// Accept() blocks the programm execution until an incoming connection request is made, then returns a net.Conn object
		// representig the connection
		connection, err := ts.listener.Accept()
		if err != nil {
			log.Print(err) // e.g. conncetion aborted
			continue
		}

		go ts.handleClient(connection) // handle connections concurrently
	}
}

func (ts *TimeServer) handleClient(connection net.Conn) {
	defer connection.Close()

	for {
		currentTimeInServerLoc := time.Now().In(ts.location).Format("Mon, 02 Jan 2006 15:04:05 MST\n")

		// Send to the client the current time with the specific format
		_, err := io.WriteString(connection, currentTimeInServerLoc)
		// If any errors occured or a client disconnected, close the connection an start waiting for a new one
		if err != nil {
			return // e.g. client disconnected
		}

		time.Sleep(time.Millisecond * 300)
	}
}
