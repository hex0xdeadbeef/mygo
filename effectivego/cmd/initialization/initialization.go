package initialization

import (
	"flag"
	"log"
	"os"
)

var (
	home   = os.Getenv("HOME")
	user   = os.Getenv("USER")
	gopath = os.Getenv("GOPATH")
)

func init() {
	if user == "" {
		log.Fatal("%USER not set")
	}

	if home == "" {
		home = "/home/" + user
	}

	if gopath == "" {
		gopath = home + "/go"
	}

	flag.StringVar(&gopath, "gopath", gopath, "override default GOPATH")
}
