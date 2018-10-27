package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
)

var ip, port *string

func init() {
	ip = flag.String("ip", "1.1.1.1", "Remote Name Server ip address")
	port = flag.String("port", "853", "Remote Port")

	flag.Parse()

	if net.ParseIP(*ip) == nil {
		fmt.Fprintln(os.Stderr, "Bad IP address")
		flag.PrintDefaults()
		os.Exit(1)
	}
	portNumber, _ := strconv.Atoi(*port)
	if !(portNumber > 0 && portNumber < 65536) {
		fmt.Fprintln(os.Stderr, "Port Number out of range")
		flag.PrintDefaults()
		os.Exit(1)
	}
}

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

	address := net.JoinHostPort(*ip, *port)
	go handleUDP(pc, address)
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err.Error())
		}
		go handleTCP(conn, address)
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
