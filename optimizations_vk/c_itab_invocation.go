package main

import "fmt"

/*
	SEARCH THE REALIZATION OF INTERFACE'S FUNCTION IN "itab"

1. Consider that there's an interface that has the only one method and there are some realizations of this in this package. In our case they are fileDumper and dbDumper. The code is light
	and easy to test.
2. The interface consist of two machine words: the pointer to the data (if the data more than 8 bytes) and the pointer to the type in the table "itable"
3. The virtual call always leads to search in "itable" and while we're searching the useful work isn't being performed at all.
4. Instead of using interfaces that lead to the search search in "itable" we can specify the enum representing kinds of objects and use it switch-case control structure to call a specific
	method. But there are some cons of it:
		1. It's not pretty-looking
		2. It's hard to expand
		3. It's hard to test
*/

type Dumper interface {
	Dump(string)
}

type fileDumper struct {
}

func (fD fileDumper) Dump(s string) {
	_ = s
}

type dbDumper struct {
}

func (dD dbDumper) Dump(s string) {
	_ = s
}

func MakeDumpsA(dumpers ...Dumper) {
	for _, v := range dumpers {
		v.Dump("hello world")
	}
}

type dumpKind int

const (
	_ dumpKind = iota
	toFile
	toDB
)

type DumperCustom struct {
	kind dumpKind
}

func DumpToDB(s string) {
	_ = s

}

func DumpToFile(s string) {
	_ = s
}

func MakeDumpsB(dumpers ...DumperCustom) {
	for _, d := range dumpers {
		switch d.kind {
		case toFile:
			DumpToFile("hello world")
		case toDB:
			DumpToDB("hello world")
		default:
			fmt.Println("invalid dumper kind")
		}
	}
}
