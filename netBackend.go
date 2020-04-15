package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

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

			var DPAK DAT
			json.Unmarshal(packetIn.Content, &DPAK)
			fmt.Print(DPAK)
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
			fmt.Println("recieved information \n")
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
	}

	for {
	}
}
