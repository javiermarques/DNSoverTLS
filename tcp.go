package main

import (
	"crypto/tls"
	"log"
	"net"
)

func handleTCP(c net.Conn) {
	//Connect to cloudflare
	dns, err := tls.Dial("tcp", "1.1.1.1:853", &tls.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}
	//Read from the client
	//and pass the same to the TLS
	buffer := read(c)
	dns.Write(buffer)

	//Send back the same respone
	response := read(dns)
	c.Write(response)

	dns.Close()
	c.Close()
}
