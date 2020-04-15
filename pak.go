package main

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

type P struct {
	PType   byte
	Content []byte
}

type DAT struct {
	PSize int
	PNum  int
	//	PHash []byte
	//	CHash []byte
	//	Data  []byte
}

type RSD struct {
	Num   int
	Group int
	CHash []byte
	PHash []byte
	CHIP  []byte
	DIP   []byte
}

type AEM struct {
	Num  int
	CHIP []byte
	DRIP []byte
	Lock []byte
}

type AEC struct {
	Num  int
	CHIP []byte
	NIP  []byte
	Lock []byte
}

type INF struct {
	IType byte
	Data  []byte
}

type INQ struct {
}

type MPQ struct {
}

type MPR struct {
}
