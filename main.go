package main

import (
	"flag"
	"log"
)

func main() {
	log.SetFlags(0)

	srvBool := flag.Bool("srv", false, "Ping Minecraft Server")
	usrBool := flag.Bool("user", false, "Get minecraft profile")

	flag.Parse()

	if *srvBool {
		McPing()
	}

	if *usrBool {
		User()
	}
}
