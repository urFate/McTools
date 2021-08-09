package main

import (
	"flag"
	"log"
)

func main() {
	log.SetFlags(0)

	srvBool := flag.Bool("srv", false, "Ping Minecraft Server (mctools -srv hypixel.net)")
	usrBool := flag.Bool("user", false, "Get Minecraft Profile (mctools -user Dinnerbone)")

	flag.Parse()

	if *srvBool {
		McPing()
	} else if *usrBool {
		User()
	} else {
		flag.Usage()
	}
}
