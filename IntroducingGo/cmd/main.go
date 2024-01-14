package main

import (
	"goprs/pkg/chapter6"
)

func main() {
	// chapter1.Hello_Go()

	// chapter2.All_Methods_Using()

	// chapter3.All_Methods_Using()

	// chapter4.All_Methods_Using()

	// chapter5.All_Methods_Using()

	// chapter6.All_Methods_Using()

	// chapter7.All_Methods_Using()

	// chapter8_9.String_All_Nethods_Using()
	// chapter8_9.Bytes_All_Nethods_Using()
	// chapter8_9.IO_All_Nethods_Using()
	// chapter8_9.List_All_Nethods_Using()
	// chapter8_9.Hash_All_Nethods_Using()
	// chapter8_9.TCP()
	// chapter8_9.HTTP()
	// chapter8_9.RPC()
	// chapter8_9.ParsingCommandLineArgs()
	// chapter10.Firs_Using_Goroutines()
	// chapter10.Channel_Using()
	// chapter10.Select_Using()
	// chapter10.Timeout_Select_Using()
	// chapter10.Buffered_Channel_Using()

	closurePointer := chapter6.MakeEvenGenerator()
	closurePointer()
	closurePointer()
}
