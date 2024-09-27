package chapter7

import (
	"bytes"
	"fmt"
	c6 "golang/pkg/chapters/chapter6"
	"io"
	"os"
	"time"
)

func InterfaceSatisfaction() {
	var w io.Writer

	w = os.Stdout
	w = new(bytes.Buffer)
	/* w = time.Second  // cannot use time.Second (constant 1000000000 of type time.Duration) as io.Writer value in
	assignment: time.Duration does not implement io.Writer (missing method Write) */

	var rwc io.ReadWriteCloser
	rwc = os.Stdout
	/* rwc = new(bytes.Buffer) // cannot use new(bytes.Buffer) (value of type *bytes.Buffer) as io.ReadWriteCloser
	value in assignment: *bytes.Buffer does not implement io.ReadWriteCloser (missing method Close) */

	fmt.Println(w, rwc)

	w = rwc // ReadWriteCloser has Write method, so everything is okay.
	/* rwc = w // cannot use w (variable of type io.Writer) as io.ReadWriteCloser value in assignment: io.Writer
	does not implement io.ReadWriteCloser (missing method Close) */

}

func StringerInterfaceExample() {
	var bitSet c6.IntSet
	bitSet.Add(1, 2, 3)

	var stringer fmt.Stringer

	/* stringer = bitSet // cannot use bitSet (variable of type chapter6.IntSet) as fmt.Stringer value in
	assignment: chapter6.IntSet does not implement fmt.Stringer (method String has pointer receiver) */
	stringer = &bitSet

	fmt.Println(stringer.String())
}

func EnvelopeExample() {
	os.Stdout.Write([]byte("Hello!"))
	os.Stdout.Close() // os.Stdout has the method Close()

	var w io.Writer
	w = os.Stdout // Now os.Stdout is assigned to variable "w", w has all the data of os.Stdout but we can call only
	// the interface's io.Writer methods
	w.Write([]byte("Hello!"))
	/* w.Close // w.Close undefined (type io.Writer has no field or method Close) */
}

func EmptyInterfaceVariable() {
	var any interface{} // interface{} has no methods it demands from implementers, so it can store all the data of
	// any type, but cannot call any methods
	any = true
	any = 12.34
	any = "hello"
	any = new(bytes.Buffer)
	any = map[string]int{"one": 1}

	bitSet := c6.IntSet{}
	bitSet.Add(1, 65, 129, 257, 513)
	any = bitSet

	fmt.Println(any)

}

/*
Assume that we have these types:
Book
Magazine

Movie
TVEpisode

Track
Podcast
*/

/*
This is the common interface:
*/

type Artifact interface {
	Title() string
	Creators() []string
	Created() time.Time
}

/*
"Personal" interfaces for each kind.
*/

// For Book/Magazine:
type Text interface {
	Pages() int
	Words() int
	PageSize() int
}

// For Track/Podcast:
type Audio interface {
	Stream() (io.ReadCloser, error)
	RunningTime() time.Duration
	Format() string
}

// For Movie/TVEpisode:
type Video interface {
	Stream() (io.ReadCloser, error)
	RunningTime() time.Duration
	Format() string
	Resolution() (x, y int)
}

/*
The union of Audio and Movie:
*/
type Streamer interface {
	Stream() (io.ReadCloser, error)
	RunningTime() time.Duration
	Format() string
}
