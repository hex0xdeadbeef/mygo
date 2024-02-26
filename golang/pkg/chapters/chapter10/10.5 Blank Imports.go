package chapter10

import (
	"flag"
	"fmt"
	"io"
	"os"

	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
)

func ConvertImage() {
	if err := toJPEG(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "jpeg: %v\n", err)
		os.Exit(1)
	}
}

func toJPEG(in io.Reader, out io.Writer) error {
	var (
		format = flag.String("f", "jpeg", "format to be applied while conversion\navailable formats: jpeg, gif, png")
	)

	flag.Parse()

	img, kind, err := image.Decode(in)
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, "Input format =", kind)

	switch *format {
	case "jpeg":
		return jpeg.Encode(out, img, &jpeg.Options{Quality: 95})
	case "png":
		return png.Encode(out, img)
	case "gif":
		return gif.Encode(out, img, &gif.Options{})
	default:
		return fmt.Errorf("unsupported format")
	}
}
