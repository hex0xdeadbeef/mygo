package mulitiersorttable

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	table *TableWidget
)

func filling() (*TableWidget, error) {
	table := &TableWidget{}
	if err := table.Init([][]string{
		{"Firstname", "Lastname", "Year", "Role"},
		{"Rustam", "Rakhmatullov", "2003", "Money lover"},
		{"Denis", "Nemchenko", "1980", "WorldOfTanks lover", "Thirsty"},
		{"Kate", "Balashova", "2003", "Noodles with carrot lover"},
		{"Dmitriy", "Mamykin", "-1000", "Triceratops", "Hungry"},
		{"George", "Odin", "2002", "Barni"},
		{"Alexander", "Zherenovskiy", "2003", "Kuni lover"},
	}); err != nil {
		return nil, fmt.Errorf("creating table; %s", err)
	}

	return table, nil
}

func StartServer() {

	var err error
	// stuff the global variable "table" with data
	table, err = filling()
	if err != nil {
		log.Printf("filling table; %s", err)
		return
	}

	log.Printf("initial table has been created")

	http.HandleFunc("/", startPage)
	http.HandleFunc("/sort", sortTablePage)
	http.ListenAndServe("localhost:8080", nil)

}

/*
sortHandle(responseWriter http.ResponseWriter, request *http.Request) handles a click on a table header an initializes sorting the table
rows
*/
func startPage(responseWriter http.ResponseWriter, request *http.Request) {
	// Send to a client table file
	sendTableFile(responseWriter, request)
}

/*
sortTablePage(responseWriter http.ResponseWriter, request *http.Request) processes a client's request in the following way:
1) Parses columnIndex parameter
2) Sets a sort parameter of the table if the parsing is successful
3) Sorts the table by the parameter has been set
4) Sends to a client the sorted table file
*/
func sortTablePage(responseWriter http.ResponseWriter, request *http.Request) {
	// Parse column index from the url parameters
	clickedColumnIndexStr := request.URL.Query().Get("columnIndex")
	clickenColumnIndexNumber, err := strconv.ParseInt(clickedColumnIndexStr, 10, 64)
	if err != nil {
		log.Printf("while parsing columnIndex; invalid columnIndex parameter: %s", clickedColumnIndexStr)
	}

	// Try to set the current table sort parameter
	if err = table.setSortParameter(int(clickenColumnIndexNumber)); err != nil {
		log.Print(err.Error())

	} else {
		// Sort table
		table.Sort()
	}

	// Send to a client the sorted table file
	sendTableFile(responseWriter, request)
}

/*
sendTableFile(responseWriter http.ResponseWriter, request *http.Request) sends to a client the a table file in a current state
*/
func sendTableFile(responseWriter http.ResponseWriter, request *http.Request) {
	// Create and fill the file with table has created above
	tableFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Printf("creating file \"%s\"; %s", fileName, err)
	}

	// Parse and fill the template with the data
	err = tableTemplate.Execute(tableFile, table)
	if err != nil {
		log.Printf("parsing and writing the filled template; %s", err)
	}
	tableFile.Close()

	// // Set status code
	// responseWriter.WriteHeader(http.StatusOK)

	// Send to the client the filled file
	http.ServeFile(responseWriter, request, fileName)
	log.Printf("file with filled table has been sent\n\n")
}
