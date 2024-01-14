package chapter8_9

import (
	"crypto/sha1"
	"fmt"
	"hash/crc32"
	"io"
	"os"
)

// "Hashes and Cryptography": hash/crc32 package
func crc32Using() {
	hasher := crc32.NewIEEE()
	hasher.Write([]byte("test"))
	sum := hasher.Sum32()
	fmt.Println("Hash value of the string \"test\":", sum)
}

// "Hashes and Cryptography": hash/crc32 package
func sha1Using() {
	hasher := sha1.New()
	hasher.Write([]byte("hello"))
	sum := hasher.Sum([]byte{})
	fmt.Println("Hash value of the string \"test\":", sum)
}

func getFileHash(fileName string) (uint32, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return 0, nil
	}
	defer file.Close()

	hash := crc32.NewIEEE()

	_, err = io.Copy(hash, file)
	if err != nil {
		return 0, nil
	}

	return hash.Sum32(), nil
}

func comparingFilesHashes(firstText, secondText string) {
	Create_WriteString("text1.txt", firstText)
	Create_WriteString("text2.txt", secondText)

	hash1, err1 := getFileHash("text1.txt")
	if err1 != nil {
		return
	}

	hash2, err2 := getFileHash("text1.txt")
	if err2 != nil {
		return
	}

	if hash1 == hash2 {
		fmt.Println("Hashes are equal")
	} else {
		fmt.Println("Hashes are different")
	}

}

func Hash_All_Nethods_Using() {

	// "Hashes and Cryptography": hash/crc32 package
	fmt.Println("\"Hashes and Cryptography\": hash/crc32 package")
	fmt.Println()

	crc32Using()
	fmt.Println()

	sha1Using()
	fmt.Println()

	comparingFilesHashes("bebra", "bebra")
	fmt.Println()

}
