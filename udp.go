package main

import (
	"crypto/tls"
	"encoding/binary"
	"log"
	"net"
)

func handleUDP(pc net.PacketConn) {
	for {
		//simple read
		buffer := make([]byte, 4096)
		n, addr, err := pc.ReadFrom(buffer)
		if err != nil {
			log.Fatal(err.Error())
		}

		buffer = buffer[:n]
		request := make([]byte, 2, n+2)
		//Adding the packet length needed for tcp
		binary.BigEndian.PutUint16(request, uint16(n))
		request = append(request, buffer...)

		dns, err := tls.Dial("tcp", "1.1.1.1:853", &tls.Config{})
		if err != nil {
			log.Fatal(err.Error())
		}
		dns.Write(request)
		response := read(dns)

		//Removing the first two bytes (size of TCP packet)
		n, err = pc.WriteTo(response[2:], addr)
		dns.Close()
	}
}
