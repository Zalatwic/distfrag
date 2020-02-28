package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

type P struct {
	PType byte
}

//spams packet type 4
func main() {
	for {
		con, err := net.Dial("tcp", "localhost:5831")
		if err != nil {
			fmt.Println(err)
		}

		packet := P{4}

		pakEncode := gob.NewEncoder(con)
		err = pakEncode.Encode(packet)

		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(2 * time.Second)

		fmt.Println("connecting")
		con.Close()
	}
}
