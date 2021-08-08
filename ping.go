package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	mcnet "github.com/Tnze/go-mc/net"
	pk "github.com/Tnze/go-mc/net/packet"
)

var (
	ip       string
	country  string
	city     string
	timezone string
	org      string
	addrType string
)

func McPing() {
	host := lookupMC(flag.Arg(0))

	conn, err := net.Dial("tcp", host[0])
	if err != nil {
		log.Default().Fatalf("[Error] Cannot connect to the server \"%v\"", flag.Arg(0))
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

		ip, addrType, country, city, timezone, org = IpInfo(addr[0].String())
	} else if !strings.Contains(flag.Arg(0), ":") {
		ip, addrType, country, city, timezone, org = IpInfo(flag.Arg(0))
	} else {
		ip, addrType, country, city, timezone, org = IpInfo(strings.Split(flag.Arg(0), ":")[0])
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

	fmt.Printf("┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  %v/%v %vms\n┃ %v\n%v\n┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n",
		status.Players.Online, status.Players.Max, delay.Milliseconds(), descriptionLine1, descriptionLine2)

	fmt.Printf("IP: %v\nType: %v\nCountry: %v\nCity: %v\nTime zone: %v\nOrg: %v\n",
		ip, addrType, country, city, timezone, org)
}

// PingAndListConn is the version of PingAndList using a exist connection.
func PingAndListConn(conn net.Conn, protocol int) (*Status, time.Duration, error) {
	addr := conn.RemoteAddr().String()
	mcConn := mcnet.WrapConn(conn)
	return pingAndList(addr, mcConn, protocol)
}

func pingAndList(addr string, conn *mcnet.Conn, protocol int) (*Status, time.Duration, error) {
	// parse hostname and port
	host, strPort, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, 0, fmt.Errorf("could not split host and port: %v", err)
	}

	port, err := strconv.ParseUint(strPort, 10, 16)
	if err != nil {
		return nil, 0, fmt.Errorf("port must be a number: %v", err)
	}

	// handshake
	err = conn.WritePacket(pk.Marshal(0x00, // packet ID
		pk.VarInt(protocol),    // protocol version
		pk.String(host),        // server host
		pk.UnsignedShort(port), // server port
		pk.Byte(1),             // next: ping
	))
	if err != nil {
		return nil, 0, fmt.Errorf("sending handshake: %v", err)
	}

	// list
	err = conn.WritePacket(pk.Marshal(0))
	if err != nil {
		return nil, 0, fmt.Errorf("sending list: %v", err)
	}

	// response
	var recv pk.Packet
	err = conn.ReadPacket(&recv)
	if err != nil {
		return nil, 0, fmt.Errorf("receiving response: %v", err)
	}

	var s pk.String
	if err = recv.Scan(&s); err != nil {
		return nil, 0, fmt.Errorf("scanning list: %v", err)
	}

	// ping
	startTime := time.Now()
	unixStartTime := pk.Long(startTime.Unix())

	err = conn.WritePacket(pk.Marshal(0x01, unixStartTime))
	if err != nil {
		return nil, 0, fmt.Errorf("sending ping: %v", err)
	}

	err = conn.ReadPacket(&recv)
	if err != nil {
		return nil, 0, fmt.Errorf("receiving pong: %v", err)
	}
	delay := time.Since(startTime)

	var t pk.Long
	if err = recv.Scan(&t); err != nil {
		return nil, 0, fmt.Errorf("scanning pong: %v", err)
	}
	// check time
	if t != unixStartTime {
		return nil, 0, errors.New("mismatched pong")
	}

	// parse status
	status := new(Status)
	if err = json.Unmarshal([]byte(s), status); err != nil {
		return nil, 0, fmt.Errorf("unmarshal json fail: %v", err)
	}

	return status, delay, nil
}

func lookupMC(addr string) (addrs []string) {
	if !strings.Contains(addr, ":") {
		_, addrsSRV, err := net.LookupSRV("minecraft", "tcp", addr)
		if err == nil && len(addrsSRV) > 0 {
			for _, addrSRV := range addrsSRV {
				addrs = append(addrs, net.JoinHostPort(addrSRV.Target, strconv.Itoa(int(addrSRV.Port))))
			}
			return
		}
		return []string{net.JoinHostPort(addr, "25565")}
	}
	return []string{addr}
}

func IpInfo(address string) (ip string, addrType string, country string,
	city string, timezone string, org string) {

	spaceClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, "https://ipwhois.app/json/"+address, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	data := IpData{}
	jsonErr := json.Unmarshal(body, &data)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return data.IP, data.Type, data.Country, data.City, data.Timezone, data.Org
}
