package main

import (
	"fmt"
	"net"
)

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	if err != nil {
		fmt.Println("Failed to resolve UDP address:", err)
		return
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Failed to bind to address:", err)
		return
	}
	defer udpConn.Close()

	buf := make([]byte, 512)

	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}
		fmt.Printf("Received %d bytes from %s\n", size, source)

		received := parse(buf[:size])

		var rcode uint8
		if received.headers.opcode != 0 {
			rcode = 4
		}
		response := NewMessage(
			&Headers{
				id:      received.headers.id,
				qr:      true,
				opcode:  received.headers.opcode,
				aa:      false,
				tc:      false,
				rd:      received.headers.rd,
				ra:      false,
				z:       0,
				rcode:   rcode,
				qdcount: 0,
				ancount: 0,
				nscount: 0,
				arcount: 0,
			})

		response.addQuestion("codecrafters.io", 1, 1)

		response.addAnswer("codecrafters.io", 1, 1, 60, [4]uint8{8, 8, 8, 8})

		_, err = udpConn.WriteToUDP(response.bytes(), source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
