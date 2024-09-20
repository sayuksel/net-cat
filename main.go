package main

import (
	"fmt"
	"netcat/pkg"
	"os"
)

func main() {
	server := &pkg.Server{}
	err := server.Create("localhost:3000", 10)
	if err != nil {
		fmt.Println("error man", err)
		os.Exit(1)
	}
	// close signale
	// go close
	for {
		conn, err := server.Server.Accept()
		if err != nil {
			fmt.Println(err)
			break
		}
		go server.Handle(conn)

	}
}
