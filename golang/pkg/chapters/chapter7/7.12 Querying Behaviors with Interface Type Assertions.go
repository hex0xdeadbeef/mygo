package chapter7

import "io"

func writeHeader(w io.Writer, contentType string) error {

	if _, err := writeString(w, "Content-Type: "); err != nil {
		return err
	}

	if _, err := writeString(w, contentType); err != nil {
		return err
	}

	// ... do some work
	return nil
}

func writeString(w io.Writer, s string) (n int, err error) {

	// The interface to check the type assertion result and apply the WriteString if the check is passed.
	type StringWriter interface {
		io.Writer
		WriteString(string) (n int, err error)
	}

	// If type is asserted write the string directly
	if multiWriter, ok := w.(StringWriter); ok {
		return multiWriter.WriteString(s)
	}

	// If type isn't asserted, allocate the new mempool, make a copy of the string data and put it into the slice.
	return w.Write([]byte(s))
}
