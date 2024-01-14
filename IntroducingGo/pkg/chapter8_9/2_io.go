package chapter8_9

import (
	"fmt"
	"os"
	"path/filepath"
)

// "os" package (Files and Folders)
func Open_Close_Stat_Read() {

	file, err := os.Open("test.txt") // Open the file and write error value into err
	if err != nil {                  // If err has unlike nil value, close the file and terminate the execution using the deferred method
		return
	}
	defer file.Close()

	stat, err := file.Stat() // Read the file info and write the value of error into err
	if err != nil {          // If value isn't nil, close the file using the deferred method and terminate the execution of program
		return
	}

	bs := make([]byte, stat.Size()) /* If all is alright, make a slice of bytes that have sufficient amount of capacity
	to store the file data */
	_, err = file.Read(bs)
	if err != nil {
		return
	}

	str := string(bs)
	fmt.Println("File stat: ", stat)
	fmt.Println("File data: ", str)
}

func ReadFile() {
	bs, err := os.ReadFile("test.txt") // There is deprecated method using
	if err != nil {
		return
	}

	str := string(bs)
	fmt.Println("File data:", str)
}

func Create_WriteString(name, text string) {
	myFirstFile, err := os.Create(name) // If everything is okay, file will be created and allocated in the working directory
	if err != nil {
		return
	}
	defer myFirstFile.Close()

	myFirstFile.WriteString(text)
}

func Readdir() {
	dir, err := os.Open("Desktop/goprs")
	if err != nil {
		return
	}
	defer dir.Close()

	listOfDirectory, err := dir.Readdir(-1) // if parameter <= 0, Readdir returns all fileInfo data
	if err != nil {
		return
	}

	for _, fileInfo := range listOfDirectory {
		fmt.Println(fileInfo.Name())
	}

}

func Walk(rootPath string) {
	filepath.Walk(rootPath,
		func(path string, info os.FileInfo, err error) error {
			fmt.Println(path)
			return nil
		})
}

func IO_All_Nethods_Using() {
	// Input/Output
	/*
		The two main interfaces "io" consists of are Reader/Writer. Reader supports reading via the Read() method and Writer
		supports sriting via the Write() method. Also "io" package has a Copy() function which copies some data from src Reader
		to dst Writer
		func Copy(dst Writer, src Reader) (written int64, err error)

	*/

	fmt.Println("--- os package (Files and Folders)")
	fmt.Println()

	fmt.Println("--- Open_Close_Stat_Read")
	Open_Close_Stat_Read()
	fmt.Println()

	fmt.Println("--- Create_WriteString")
	Create_WriteString("firstFile.txt", "This is my first file.")
	fmt.Println()

	fmt.Println("--- Readdir")
	Readdir()
	fmt.Println()

	fmt.Println("--- Walk")
	Walk("/Users/dmitriymamykin/Library/Mobile Documents/com~apple~CloudDocs/Go/goprojects")
	fmt.Println()

}
