package logs

import (
	"log"
	"os"
)

func panicFatal() {
	if len(os.Args) == 1 {
		log.Fatal("Fatal")
	}

	log.Panic("Panic")
}

func Using(vals ...string) {
	os.Args = append(os.Args, vals...)

	panicFatal()
}
