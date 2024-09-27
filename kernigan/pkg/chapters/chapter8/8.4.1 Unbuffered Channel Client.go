package chapter8

import (
	"io"
	"log"
	"net"
	"os"
)

func StartUnbufferedChannelConnection() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan struct{})

	go func() {
		_, err := io.Copy(os.Stdout, conn) // Ignoring errors
		if err != nil {
			log.Println(err)
		}

		log.Println("done")

		done <- struct{}{}

	}()

	mustCopy(conn, os.Stdin)
	// conn.Close()
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.CloseWrite()
	}

	<-done
}
