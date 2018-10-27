package main

import (
	"crypto/tls"
	"io"
	"log"
	"net"
)

func main() {
	l, err := net.Listen("tcp", ":53")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer l.Close()

	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err.Error())
		}
		go func(c net.Conn) {
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
		}(conn)
	}
}

func read(c net.Conn) []byte {
	//big buffer
	response := make([]byte, 0, 4096)
	//temporal buffer to read over
	tmp := make([]byte, 4)
	for {
		n, err := c.Read(tmp)
		if err != nil {
			if err != io.EOF {
				log.Fatal(err.Error())
			}
			break
		}
		response = append(response, tmp[:n]...)
		//net.Conn doesn't send EOF
		//so we have to break when read is completed
		if n < 4 {
			break
		}

	}
	return response
}
