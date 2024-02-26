package zipreader

import (
	"archive/zip"
	"fmt"
	"golang/pkg/chapters/chapter10/archivereader"
	"io"
)

func init() {
	archivereader.Register("zip", NewZipReader)
}

type ZipReader struct {
	reader *zip.ReadCloser
}

func (zr *ZipReader) Open(path string) error {
	zipReader, err := zip.OpenReader(path)
	if err != nil {
		return err
	}

	zr.reader = zipReader
	return nil
}

func (zr *ZipReader) Read(writer io.Writer) error {
	const (
		tableHeaderFormat = "%-*s %-*s\n"
		tableRowFormat    = "%-*s %-*d\n"
		width             = 40
	)

	fmt.Fprintf(writer, tableHeaderFormat, width, "Name", width, "Size")

	for _, file := range zr.reader.File {
		fmt.Fprintf(writer, tableRowFormat, width, file.Name, width, file.UncompressedSize64)
	}
	return nil
}

func (zr *ZipReader) Close() error {
	return zr.reader.Close()
}

func NewZipReader() archivereader.ArchiveReader {
	zipReader := &ZipReader{}
	return zipReader
}

/*
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
*/
