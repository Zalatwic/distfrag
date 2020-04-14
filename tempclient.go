package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

//spams packet type 0
func main() {
	for {
		con, err := net.Dial("tcp", "localhost:5831")
		if err != nil {
			fmt.Println(err)
		}

		dat := DAT{
			PSize: 200,
			PNum:  40,
		}

		hi, _ := json.Marshal(dat)

		packet := P{0, hi}

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
