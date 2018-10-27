package main

import (
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

	// listen to incoming udp packets
	pc, err := net.ListenPacket("udp", ":53")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer pc.Close()

	go handleUDP(pc)
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err.Error())
		}
		go handleTCP(conn)
	}
}

//helper function to read tcp connection
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
