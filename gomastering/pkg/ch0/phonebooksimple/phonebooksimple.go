package phonebooksimple

import (
	"fmt"
	"os"
	"path"
)

type entry struct {
	name    string
	surname string
	tel     string
}

var (
	entries = []entry{}
)

func search(key string) *entry {
	for i, e := range entries {
		if e.surname == key {
			return &entries[i]
		}
	}

	return nil
}

func list() {
	for _, e := range entries {
		fmt.Println(e)
	}
}

func Using() {
	args := os.Args

	if len(args) == 1 {
		exe := path.Base(args[0])
		fmt.Printf("Usage %s search|list <args>\n", exe)
		return
	}

	entries = append(entries, entry{name: "Dmitriy", surname: "Mamykin", tel: "1"})
	entries = append(entries, entry{name: "Rustam", surname: "Rakhmatullov", tel: "2"})

	switch args[1] {
	case "search":
		if len(args) != 3 {
			fmt.Println("Usage: search Surname")
			return
		}

		if result := search(args[2]); result != nil {
			fmt.Println("Entry not found:", args[2])
			return
		} else {
			fmt.Println(*result)
		}
	case "list":
		list()

	default:
		fmt.Println("Not a valid option")
	}
}
