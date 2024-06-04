package interfacenames

type CustomWriterReaderCloser interface {
	Write(b []byte) (int, error)
	Read(b []byte) (int, error)
	Close() error
}

type StringConverter interface {
	String() string
}
