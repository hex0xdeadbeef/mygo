package chatsnippet

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

var (
	cl CustomLogger
)

func init() {
	cl = *getLogger()
}

type client struct {
	name string
	conn net.Conn
	c    chan string
}

func (client *client) setName() {
	fmt.Fprint(client.conn, "\nENTER YOUR NAME: ")
	scanner := bufio.NewScanner(client.conn)
	for {
		if scanner.Scan() {
			client.name = scanner.Text()
			return
		} else if err := scanner.Err(); err != nil {
			log.Printf("while client with address %s: %s", client.conn.LocalAddr(), err)
		}
	}
}

func (client *client) getMessage() {
	for msg := range client.c {
		fmt.Fprintln(client.conn, msg)
	}
}

func (client *client) closeActions() {
	client.conn.Close()
	close(client.c)
}

var (
	wg   sync.WaitGroup
	done = make(chan struct{})

	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
)

func ChatServer() {
	
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		broadcaster()
	}()

	wg.Add(1)
	go func() {
		wg.Done()
		<-done
		listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			break
		}

		wg.Add(1)
		go func() {
			wg.Done()
			handleConn(conn)
		}()
	}

	wg.Wait()
}

func broadcaster() {
	const (
		serverlTimeoutInMinutes = 1
		sendSleepTimeInMillis   = 50
	)

	var (
		clients = make(map[client]bool)
		ticker  = time.NewTicker(serverlTimeoutInMinutes * time.Minute)
	)

	cl.logger.Println(logSplitter)
	cl.logger.Println(sessionStartedMsg)

	defer func() {
		ticker.Stop()

		if err := cl.Close(); err != nil {
			panic(err)
		}
	}()

outer:
	for {
		select {
		case newClient := <-entering:
			newClient.c <- announceAllClients(clients)
			clients[newClient] = true

		case leftClient := <-leaving:
			delete(clients, leftClient)

			leftClient.closeActions()

		case msg := <-messages:
			cl.logger.Println(msg)
			for client := range clients {
			inner:
				for {
					select {
					case client.c <- msg:
						break inner
					default:
						time.Sleep(sendSleepTimeInMillis * time.Millisecond)
					}
				}
			}

		case <-ticker.C:
			close(done)
			break outer
		}

	}

	cl.logger.Println(sessionCloseddMsg)
	cl.logger.Println(logSplitter)
	cl.logger.Println()
}

func handleConn(conn net.Conn) {
	const (
		clientBufferSize = 8

		tickerDurationInSeconds = 1
		timeoutInSeconds        = 30
	)

	var (
		client = client{conn: conn, c: make(chan string, clientBufferSize)}

		// Utilities for sending timeout
		ticker      = time.NewTicker(tickerDurationInSeconds * time.Second)
		messageSent = make(chan struct{})
	)

	client.setName()

	// Utilities closures
	defer func() {
		ticker.Stop()
		close(messageSent)
	}()

	messages <- fmt.Sprintf(welcomeMsgTemplate, client.name)

	go func() {
		client.getMessage()
	}()

	client.c <- fmt.Sprintf("YOU ARE: %s", client.name)
	entering <- client

	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			messages <- fmt.Sprintf("%s: %s", client.name, scanner.Text())
			messageSent <- struct{}{}
		}
	}()

	for i := timeoutInSeconds; i > 0; i-- {
		select {
		case <-ticker.C:
		case <-messageSent:
			i = timeoutInSeconds
		}
	}
	ticker.Stop()

	messages <- fmt.Sprintf(goodbyeMsgTemplate, client.name)
	leaving <- client
}

func announceAllClients(clients map[client]bool) string {
	var (
		currentClients bytes.Buffer
	)

	currentClients.WriteString("USERS ONLINE:\n")

	for client := range clients {
		currentClients.WriteString(client.name)
		currentClients.WriteByte('\n')
	}

	return currentClients.String()
}
