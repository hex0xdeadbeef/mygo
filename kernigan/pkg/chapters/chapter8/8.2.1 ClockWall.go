package chapter8

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

/*
cd /Users/dmitriymamykin/Desktop/client

1) go run main.go NewYork=localhost:8010 London=localhost:8020 Tokyo=localhost:8030

2) go run main.go NewYork=localhost:8010 Los_Angeles=localhost:8040 Sydney=localhost:8050

3) go run main.go NewYork=localhost:8010 Tokyo=localhost:8040
*/

const (
	splitSet           = "=localhost:"
	minValidArgsCount  = 2
	validPortionsCount = 2
	duration           = 300
)

type ClockWall struct {
	zones []string
	ports []string

	connections []net.Conn
}

func NewClockWall() (*ClockWall, error) {
	if len(os.Args) < minValidArgsCount {
		return nil, nil
	}

	newClockWall := ClockWall{}
	for _, arg := range os.Args[1:] {
		argPortions, err := isLocationPortParamValid(arg)
		if err != nil {
			return nil, fmt.Errorf("validating argument; %s", err)
		}
		newClockWall.zones, newClockWall.ports =
			append(newClockWall.zones, argPortions[0]), append(newClockWall.ports, argPortions[1])
	}

	return &newClockWall, nil
}

func isLocationPortParamValid(arg string) ([]string, error) {
	portions := strings.Split(arg, splitSet)

	if portions[0] == arg {
		return nil, fmt.Errorf("invalid argument; %s", arg)
	}

	if len(portions) != validPortionsCount {
		return nil, fmt.Errorf("invalid argument; %s", arg)
	}

	_, err := parsePortNumber(portions[1])
	if err != nil {
		return nil, fmt.Errorf("invalid port number; %s", portions[1])
	}

	return portions, nil
}

func (cw *ClockWall) GetTable() {
	cw.connections = make([]net.Conn, len(cw.ports))

	for i, port := range cw.ports {
		connection, err := net.Dial("tcp", serverPrefix+port)
		if err != nil {
			log.Print(err)
			continue
		}
		defer connection.Close()

		cw.connections[i] = connection
	}

	for {
		cw.PrintTable()
		time.Sleep(time.Second * 1)
	}

}

func (cw *ClockWall) PrintTable() {

	const (
		width int = 20
	)
	var (
		table   string
		border  = strings.Repeat("-", 50)
		newLine = "\n"
	)

	for {
		table = ""
		table += border + newLine
		for i := 0; i < len(cw.connections); i++ {

			if cw.connections[i] == nil {
				continue
			}
			table += fmt.Sprintf("%-*s | ", width, cw.zones[i])

			tableRow, err := getLocalTime(cw.connections[i])
			if err != nil {
				tableRow = "-"
			}

			table += fmt.Sprintf("%-*s", width, tableRow)

		}
		table += border + newLine

		fmt.Println(table)

		time.Sleep(time.Second * 3)
	}

}

func getLocalTime(connection net.Conn) (string, error) {
	reader := bufio.NewReader(connection)

	data, err := reader.ReadBytes('\n')
	if err != nil {
		return "-", err
	}

	return string(data), err
}
