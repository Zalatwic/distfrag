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
	"io/ioutil"
	"math"
	"os"
)

type Bank struct {
	PakLen int
}

type OutAbout struct {
	Key           string
	Nonce         string
	TrailingBytes int
	FileName      string
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

//load the key file into the bank
func loadKey(file string) OutAbout {
	var rBank OutAbout
	tBank, _ := loadFile(file)

	json.Unmarshal(tBank, &rBank)

	return rBank
}

//generate hex string for aes256
func genHex(length int) string {
	oByte := make([]byte, length)
	_, err := rand.Read(oByte)
	errorCheck(err)

	return hex.EncodeToString(oByte)
}

//generate encoded packets, save a json and info
func genEncodedPackets(iFile string, pSize int) (string, [][]byte) {
	keyStr := genHex(32)
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

	nameHex := genHex(8)

	outTracker := OutAbout{
		Key:           keyStr,
		Nonce:         hex.EncodeToString(nonce),
		TrailingBytes: remainingPak,
		FileName:      iFile,
	}

	var jsonOut []byte
	jsonOut, err = json.Marshal(outTracker)
	errorCheck(err)

	err = ioutil.WriteFile("KEYS/"+nameHex+".key", jsonOut, 0644)
	errorCheck(err)

	fmt.Printf("%x\n", eDat[0])
	return nameHex, eDat
}

//take encoded packets and a keyfile then save a file
func genDecodedFile(keyName string, pak [][]byte) {
	keyRing := loadKey("KEYS/" + keyName + ".key")
	netInf := loadConfig("netConf.json")
	fmt.Println(len(pak))
	outFile := make([]byte, 0)
	rawPak := make([]byte, netInf.PakLen)

	key, _ := hex.DecodeString(keyRing.Key)
	nonce, _ := hex.DecodeString(keyRing.Nonce)

	block, err := aes.NewCipher(key)
	errorCheck(err)
	gcmCore, err := cipher.NewGCM(block)
	errorCheck(err)

	for i := 0; i < len(pak); i++ {
		rawPak, _ = gcmCore.Open(nil, nonce, pak[i], nil)

		if i == len(pak)-1 {
			rawPak = append([]byte(nil), rawPak[:keyRing.TrailingBytes]...)
		}

		outFile = append(outFile, rawPak...)
	}

	ioutil.WriteFile("PROCESSED/"+keyRing.FileName, outFile, 0644)
	return
}

func main() {
	bank := loadConfig("netConf.json")
	byteFile, _ := loadFile("test.png")
	fmt.Println(bank.PakLen)
	fmt.Println(byteFile)
	genDecodedFile(genEncodedPackets("test.png", bank.PakLen))
}
