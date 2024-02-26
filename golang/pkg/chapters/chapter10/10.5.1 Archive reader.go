package chapter10

import (
	"flag"
	"fmt"
	"golang/pkg/chapters/chapter10/archivereader"
	_ "golang/pkg/chapters/chapter10/archivereader/tarreader"
	_ "golang/pkg/chapters/chapter10/archivereader/zipreader"
)

func PrintArchiveNames() error {
	var (
		archPath = flag.String("path", "", "path to archive")
		format   = flag.String("f", "zip", "available: zip/tar")
	)

	flag.Parse()

	switch *format {
	case "zip":
		return archivereader.ReadArchive("zip", *archPath)
	case "tar":
		return archivereader.ReadArchive("tar", *archPath)
	default:
		return fmt.Errorf("unsupported format")
	}

}
