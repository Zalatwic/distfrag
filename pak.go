package main

import (
	"encoding/gob"
	"fmt"
)

// IType	Information packet type
// PType	Packet type
// PSize	Packet size
// PNum		Packet number in sequence
// PHash	File hash
// CHash	Chunk hash
// Num		Cluster number
// Group	Cluster group assignment for packet
// CHIP		Cluster head IP
// DIP		Destination IP
// DRIP		Dropped IP
// NIP		New IP for so-cluster head
// Lock		AORN, no idea, but i might use it in the future

type DAT struct {
	PType byte
	PSize int
	PNum  int
	PHash []byte
	CHash []byte
	Data  []byte
}

type RSD struct {
	PType byte
	Num   int
	Group int
	CHash []byte
	PHash []byte
	CHIP  []byte
	DIP   []byte
}

type AEM struct {
	PType byte
	Num   int
	CHIP  []byte
	DRIP  []byte
	Lock  []byte
}

type AEC struct {
	PType byte
	Num   int
	CHIP  []byte
	NIP   []byte
	Lock  []byte
}

type INF struct {
	PType byte
	IType byte
	Num   int
	PNum  int
	CHIP  []byte
	Data  []byte
}

type INQ struct {
	PType byte
}

type MPQ struct {
	PType byte
}

type MPR struct {
	PType byte
}

func main() {
	fmt.Println("vim-go")
}
