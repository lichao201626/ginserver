package main

import (
	"fmt"
	"net"
	// "ipc/rpc"
)

func main() {
	fmt.Println("client")
	con, _ := net.Dial("tcp4", "127.0.0.1:3333")
	fmt.Println("aa")
	i := 0
	for i < 5 {
		con.Write([]byte("ping"))
		// defter con.Close()
		fmt.Println("connected")
		bf := make([]byte, 2000)
		x, _ := con.Read(bf)
		fmt.Println("client received", string(bf[:x]))
		i++
	}

}
