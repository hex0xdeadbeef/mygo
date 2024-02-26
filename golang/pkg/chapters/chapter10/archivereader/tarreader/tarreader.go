package tarreader

import (
	"archive/tar"
	"errors"
	"fmt"
	"golang/pkg/chapters/chapter10/archivereader"
	"io"
	"os"
)

func init() {
	archivereader.Register("tar", NewTarReader)
}

type TarReader struct {
	file   *os.File
	reader *tar.Reader
}

func (tr *TarReader) Open(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	tr.file = file

	tr.reader = tar.NewReader(file)
	return nil
}

func (tr *TarReader) Read(writer io.Writer) error {
	const (
		tableHeaderFormat = "%-*s %-*s\n"
		tableRowFormat    = "%-*s %-*d\n"
		width             = 40
	)

	fmt.Fprintf(writer, tableHeaderFormat, width, "Name", width, "Size")

	for {
		header, err := tr.reader.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			} else {
				return err
			}
		}
		fmt.Fprintf(writer, tableRowFormat, width, header.Name, width, header.Size)
	}
}

func (tr *TarReader) Close() error {
	return tr.file.Close()
}

func NewTarReader() archivereader.ArchiveReader {
	tarReader := &TarReader{}
	return tarReader
}
