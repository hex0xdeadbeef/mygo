package archivereader

import (
	"fmt"
	"io"
	"os"
)

type ArchiveReader interface {
	Open(path string) error
	Read(writer io.Writer) error
	Close() error
}

const (
	tableHeaderFormat = "%-*s%-*s\n"
	tableRowFormat    = "%-*s%-*d\n"
	width             = 15
)

var (
	archiveReaders = make(map[string]ArchiveReader)
)

func Register(format string, newReader func() ArchiveReader) {
	archiveReaders[format] = newReader()
}

func ReadArchive(format string, path string) error {
	reader, ok := archiveReaders[format]
	if !ok {
		return fmt.Errorf("unsupported format")
	}

	if err := reader.Open(path); err != nil {
		return fmt.Errorf("opening %s: %s", path, err)
	}

	if err := reader.Read(os.Stdout); err != nil {
		return fmt.Errorf("reading %s: %s", path, err)
	}

	if err := reader.Close(); err != nil {
		return fmt.Errorf("closing %s: %s", path, err)
	}

	return nil
}
