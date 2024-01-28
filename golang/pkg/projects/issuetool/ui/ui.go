package issuetool

import (
	"fmt"
	"os"
)

func AuthPrint() {
	fmt.Fprintf(os.Stdout, "\nYOU'VE BEEN SUCCESSFULLY AUTHORIZED\n")
}

func NonAuthPrint() {
	fmt.Fprintf(os.Stdout, "\nAUTHORIZATION FAILED\n")
}
