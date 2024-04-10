package chapter8_9

import (
	"encoding/gob"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
)

func TCPServer() {
	listener, err := net.Listen("tcp", ":9999") // Listener returs us a listener that implemets the Listener interface
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		// Accept a connection
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleTCPServerConnection(connection)
	}
}

func handleTCPServerConnection(connection net.Conn) {
	message := ""

	err := gob.NewDecoder(connection).Decode(&message)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Received:", message)
	}

	connection.Close()
}

func TCPclient(connectionsAmount int) {
	getConnection := func() (net.Conn, error) {
		connection, err := net.Dial("tcp", "127.0.0.1:9999")

		if err != nil {
			fmt.Println(err)
			return nil, nil
		}

		return connection, nil

	}

	for i := 0; i < connectionsAmount; i++ {
		connection, err := getConnection()
		if err != nil {
			fmt.Println(err)
		}

		var message = ""
		message += strconv.FormatFloat(rand.Float64(), 'f', -1, 64)

		fmt.Println("Sending", message)
		err = gob.NewEncoder(connection).Encode(message)
		if err != nil {
			fmt.Println(err)
		}
		connection.Close()
	}
}

func TCP() {
	go TCPServer()
	go TCPclient(1)

	input := ""
	fmt.Scanln(&input)
}

func HTTPhello(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/plain")
	io.WriteString(res, "hello")
}

func HTTP() {
	http.HandleFunc("/HTTPhello", HTTPhello)

	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

type Server struct{} // A type that has only one fuction

func (this *Server) Negate(number int64, reply *int64) error {
	*reply = -number
	return nil
}

func RPCserver() {
	rpc.Register(new(Server))
	listener, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		connection, err := listener.Accept()
		if err != nil {
			continue
		}

		go rpc.ServeConn(connection)
	}
}

func RPCclient() {
	connection, err := rpc.Dial("tcp", "127.0.0.1:9999")
	if err != nil {
		fmt.Println(err)
		return
	}

	var result int64
	err = connection.Call("Server.Negate", int64(999), &result)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Server.Negate(999)", result)
	}
}

func RPC() {
	go RPCserver()
	go RPCclient()

	input := ""
	fmt.Scanln(&input)
}
