package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	mcnet "github.com/Tnze/go-mc/net"
)

var (
	ip       string
	hostname string
	country  string
	city     string
	timezone string
	org      string
)

func main() {
	flag.Parse()
	host := lookupMC(flag.Arg(0))

	conn, err := net.Dial("tcp", host[0])
	if err != nil {
		log.Fatal("[Error] Unknown host")
	}
	status, delay, err := pingAndList(host[0], mcnet.WrapConn(conn), 758)
	if err != nil {
		log.Fatalf("[Error] %v\n", err)
	}

	if !strings.Contains(flag.Arg(0), ":") && net.ParseIP(flag.Arg(0)) == nil {
		var targetAddr string

		_, addrsSRV, err := net.LookupSRV("minecraft", "tcp", flag.Arg(0))
		if err == nil {
			targetAddr = addrsSRV[0].Target
		} else {
			addrA, err := net.LookupIP(flag.Arg(0))
			if err == nil {
				targetAddr = addrA[0].String()
			} else {
				addrCNAME, err := net.LookupCNAME(flag.Arg(0))
				if err == nil {
					targetAddr = addrCNAME
				}
			}
		}

		addr, _ := net.LookupIP(targetAddr)

		ip, hostname, country, city, timezone, org = IpInfo(addr[0].String())
	} else if !strings.Contains(flag.Arg(0), ":") {
		ip, hostname, country, city, timezone, org = IpInfo(flag.Arg(0))
	} else {
		ip, hostname, country, city, timezone, org = IpInfo(strings.Split(flag.Arg(0), ":")[0])
	}

	var description = strings.Split(status.Description.String(), "\n")
	var descriptionLine1 string
	var descriptionLine2 string

	if strings.Contains(status.Description.String(), "\n") {
		descriptionLine1 = description[0]
		descriptionLine2 = "\033[0m┃ " + description[1]
	} else {
		descriptionLine1 = description[0]
		descriptionLine2 = "\033[0m┃"
	}

	fmt.Printf("┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ %v/%v %vms\n┃ %v\n%v\n┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n",
		status.Players.Online, status.Players.Max, delay.Milliseconds(), descriptionLine1, descriptionLine2)

	fmt.Printf("Hostname: %v\nIP: %v\nCountry: %v\nCity: %v\nTime zone: %v\nOrg: %v\n",
		hostname, ip, country, city, timezone, org)
}
