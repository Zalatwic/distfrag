package main

import (
	"bufio"
	"fmt"
	"os"
)

func errorCheck(x error) {
	if x != nil {
		panic(x)
	}
}

//this loads a file from namestring and returns it as a byte slice
func loadFile(fName string) []byte {
	readDat, err := os.Open(fName)
	errorCheck(err)
	fileInfo, err := readDat.Stat()
	errorCheck(err)

	defer readDat.Close()

	var fileSize int64 = fileInfo.Size()
	bufOut := make([]byte, fileSize)
	bufIn := bufio.NewReader(readDat)
	_, err = bufIn.Read(bufOut)

	fmt.Println("Read", fileSize, "bytes with error code:", err)
	return bufOut
}

func main() {
	byteFile := loadFile("test.txt")

}
