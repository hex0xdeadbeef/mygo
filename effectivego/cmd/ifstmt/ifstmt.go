package ifstmt

import (
	"fmt"
	"io/fs"
	"os"
)

func errorLines() error {
	f, err := os.Open("file")
	if err != nil {
		return err
	}

	d, err := f.Stat()
	if err != nil {
		f.Close()
		return err
	}

	codeUsing(f, d)

	return nil
}

func codeUsing(f *os.File, d fs.FileInfo) {
	fmt.Println(f)
	fmt.Println(d)
}
