package chapter8_9

import (
	"bytes"
	"fmt"
)

func bufferAndWriteUsing(source string) {
	var buf bytes.Buffer
	lengthOfWrittenString, err := buf.Write([]byte(source))
	fmt.Println("The data buffer contains: ", buf)
	fmt.Println("Length of the written string is:", lengthOfWrittenString, "Data that is contained in err:", err)
	fmt.Println("Now we have converted the data in the buffer and got:", buf.String())
}

func Bytes_All_Nethods_Using() {

	bufferAndWriteUsing("I'm Dmitriy")
	fmt.Println()
}
