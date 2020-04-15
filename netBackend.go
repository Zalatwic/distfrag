package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type CluHead struct {
	IP  string
	Num int
}

//CH -> is cluster head
var CH = false
var subComp []string
var coHeads []CluHead
var groupCounts = []int{0, 0, 0}
var peerCounts = []int{0, 0, 0, 0, 0, 0}

var numClusters = 3

//spawned function to take care of incoming packets
func rpak(con net.Conn) {
	fmt.Printf("connected to %s\n", con.RemoteAddr().String())

	for {
		gobDecode := gob.NewDecoder(con)
		var packetIn P
		err := gobDecode.Decode(&packetIn)
		fmt.Println(packetIn)
		packetType := packetIn.PType

		if err != nil {
			fmt.Println(err)
			fmt.Println("fucl")
			return
		}

		//packets are identified by the first byte and handled accordingly

		//DATa packet
		if packetType == 0 {
			fmt.Println("recieved data packet \n")

			var PAK DAT
			json.Unmarshal(packetIn.Content, &PAK)
			fmt.Print(PAK)
			break
		}

		//Request to Send Data
		if packetType == 1 {
			fmt.Println("recieved a request to send data \n")
			break
		}

		//Action of Elevation Message
		if packetType == 2 {
			fmt.Println("recieved notice, becoming a cluster head \n")
			break
		}

		//Action of Elevation Confirmation
		if packetType == 3 {
			fmt.Println("recieved notice, requested new cluster head added to network \n")
			break
		}

		//INFo packet
		if packetType == 4 {
			var PAK INF
			json.Unmarshal(packetIn.Content, &PAK)
			fmt.Print(PAK)

			fmt.Println("recieved information \n")

			//add connected computer to peer list
			if PAK.IType == 0 {
				if CH {
					if len(coHeads) < (numClusters * 2) {
						groupAssignment := 0
						temp := 3

						for i := 0; i < numClusters; i++ {
							if temp > groupCounts[i] {
								groupAssignment = i
								temp = groupCounts[i]
							}
						}

						//add as cohead, send out elevation message, inform coheads of new addition
						coHeads = append(coHeads, string(PAK.Data))
						go anointHead(string(PAK.Data))
						go informOthersAnointed(string(PAK.Data))

					} else {
						//assign to cohead with lowest number of connected peers

					}
				} else {
					//pass information to cluster head for assignment

				}
			}
			break
		}

		//INfo packet Query (or INQuery, whatever you prefer)
		if packetType == 5 {
			fmt.Println("recieved a request for information \n")
			break
		}

		//Marco Polo Question
		if packetType == 6 {
			fmt.Println("recieved a marco-polo question \n")
			break
		}

		//Marco Polo Responce
		if packetType == 7 {
			fmt.Println("recieved a marco-polo solution \n")
			break
		}
	}
	fmt.Println("closed pipe")
	con.Close()
}

func sendPak(adrs string, pakCon []byte, pType byte) {
	con, err := net.Dial("tcp", adrs+":5831")
	errorCheck(err)

	packet := P{pType, pakCon}
	pakEncode := gob.NewEncoder(con)
	err = pakEncode.Encode(packet)
	errorCheck(err)
}

func anointHead(NHA string) {

}

func informOthersAnointed(NHA string) {

}

func acceptConnect() {
	l, err := net.Listen("tcp4", ":5831")

	if err != nil {
		fmt.Println(err)
		return
	}

	defer l.Close()

	for {
		con, err := l.Accept()

		if err != nil {
			fmt.Println(err)
			return
		}

		go rpak(con)

	}

}

func main() {
	go acceptConnect()

	if len(os.Args) > 1 {
		var infContainer []byte
		originInfPacket := INF{
			IType: 0,
			PNum:  0,
			Data:  []byte(os.Args[1]),
		}

		infContainer, _ = json.Marshal(originInfPacket)
		sendPak(os.Args[2], infContainer, 4)
	} else {
		CH = true
	}

	for {
	}
}
