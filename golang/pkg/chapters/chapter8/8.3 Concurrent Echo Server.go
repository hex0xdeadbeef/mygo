package chapter8

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

func StartEchoServer() {

	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return // connection aborted
		}

		// Pass a connection to the handler
		go handleFunc(conn)
	}
}

func handleFunc(conn net.Conn) {

	const (
		maxTimeForWaitingInSeconds = 3
		echosTime                  = 3
	)

	var (
		input = bufio.NewScanner(conn)
		shout = make(chan struct{})
		wg    sync.WaitGroup
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for input.Scan() {
			shout <- struct{}{}

			wg.Add(1)
			// echo(conn, input.Text(), 1*time.Second) // reverb1
			go func() {
				defer wg.Done()
				echo(conn, input.Text(), 1*time.Second) // reverb2
			}()
		}
	}()

	// Wait for the client's shout. If it won't sound after 10 seconds waiting, the connection will be closed.
	ticker := time.NewTicker(1 * time.Second)
	for timeLeft := maxTimeForWaitingInSeconds; timeLeft > 0; timeLeft-- {
		select {
		case <-ticker.C:
		case <-shout:
			timeLeft = maxTimeForWaitingInSeconds + echosTime
		}
	}
	ticker.Stop()
	conn.Close()

	go func() {
		wg.Wait()
		if netConn, ok := conn.(*net.TCPConn); ok {
			netConn.CloseWrite()
		}
	}()

}

func echo(conn net.Conn, shout string, delay time.Duration) {
	fmt.Fprintln(conn, "\t", strings.ToUpper(shout))
	time.Sleep(delay)

	fmt.Fprintln(conn, "\t", shout)
	time.Sleep(delay)

	fmt.Fprintln(conn, "\t", strings.ToLower(shout))
	time.Sleep(delay)
}
