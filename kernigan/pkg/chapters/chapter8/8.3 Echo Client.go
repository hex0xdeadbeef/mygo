package chapter8

import (
	"io"
	"log"
	"net"
	"os"
)

func StartEchoConnection() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Wait to data from conn
	go mustCopy(os.Stdout, conn)

	// Write the data to pass to the server
	mustCopy(conn, os.Stdin)

}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
