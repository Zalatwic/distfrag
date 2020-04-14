//using aes-gcm cipher

package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"os"
)

type Bank struct {
	PakLen int
}

type OutAbout struct {
	Key              string
	Nonce            string
	RemainingPackets int64
	FileName         string
}

func errorCheck(x error) {
	if x != nil {
		panic(x)
	}
}

//this loads a file from namestring and returns it as a byte slice
func loadFile(fName string) ([]byte, int64) {
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
	return bufOut, fileSize
}

//load the config file into the bank
func loadConfig(file string) Bank {
	var rBank Bank
	tBank, _ := loadFile(file)

	json.Unmarshal(tBank, &rBank)

	return rBank
}

func genHex() string {
	oByte := make([]byte, 32)
	_, err := rand.Read(oByte)
	errorCheck(err)

	return hex.EncodeToString(oByte)
}

func genEncodedPackets(iFile string, pSize int) [][]byte {
	keyStr := genHex()
	key, _ := hex.DecodeString(keyStr)

	fileBytes, iSize := loadFile(iFile)
	numPak := int(math.Ceil(float64(iSize) / float64(pSize)))
	eDat := [][]byte{}
	nonce := make([]byte, 12)

	_, err := rand.Read(nonce)
	errorCheck(err)

	block, err := aes.NewCipher(key)
	errorCheck(err)
	gcmCore, err := cipher.NewGCM(block)
	errorCheck(err)

	remainingPak := 0
	for i := 0; i < numPak; i++ {
		if i == numPak-1 {
			szuntPak := make([]byte, pSize)
			remainingPak = copy(szuntPak, fileBytes[(pSize*i):])
			eDat = append(eDat, gcmCore.Seal(nil, nonce, szuntPak, nil))
		} else {
			eDat = append(eDat, gcmCore.Seal(nil, nonce, fileBytes[(pSize*i):(pSize*(i+1))], nil))
		}
	}

	fmt.Printf("%x\n", eDat)
	return eDat
}

func main() {
	bank := loadConfig("netConf.json")
	byteFile, _ := loadFile("test.txt")
	fmt.Println(bank.PakLen)
	fmt.Println(byteFile)
	genEncodedPackets("test.txt", bank.PakLen)
}
