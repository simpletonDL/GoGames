package main

import (
	"github.com/simpletonDL/GoGames/server"
	"os"
)

func main() {
	args := os.Args
	if len(args) == 1 {
		println("Please provide port number (like 8080)")
		return
	}
	port := ":" + args[1]
	println("I am multi-player server! Starting on port " + port)

	server.Run(port)
}
