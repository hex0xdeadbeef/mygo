package params

import (
	"log"
	"net/http"
)

func StartServer() {

	http.HandleFunc("/search", search)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
