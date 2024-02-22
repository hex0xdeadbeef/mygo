package chapter7

import "io"

type StringReader struct {
	str   string
	index int
}

func (sr *StringReader) String() string {
	return sr.str
}

func (sr *StringReader) Read(p []byte) (n int, err error) {
	if sr.index >= len(sr.str) {
		return 0, io.EOF
	}

	n = copy(p, sr.str[sr.index:])
	sr.index += n

	return n, nil
}

func NewReader(source string) io.Reader {
	newStringReader := &StringReader{str: source, index: 0}

	return newStringReader
}

type LimitReader struct {
	ByteLimit int64
	Reader    io.Reader
}

func (lr *LimitReader) Read(p []byte) (n int, err error) {
	if lr.ByteLimit <= 0 {
		return 0, io.EOF
	}
	if int64(len(p)) > lr.ByteLimit {
		p = p[:lr.ByteLimit]
	}

	n, err = lr.Reader.Read(p)
	lr.ByteLimit -= int64(n)

	return
}

func NewLimitedReader(r io.Reader, n int64) io.Reader {
	newLimitReader := &LimitReader{Reader: r, ByteLimit: n}

	return newLimitReader
}
