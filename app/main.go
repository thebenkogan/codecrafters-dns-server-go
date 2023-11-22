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

		receivedData := string(buf[:size])
		fmt.Printf("Received %d bytes from %s: %s\n", size, source, receivedData)

		response := NewMessage(
			&Headers{
				id:      1234,
				qr:      true,
				opcode:  0,
				aa:      false,
				tc:      false,
				rd:      false,
				ra:      false,
				z:       0,
				rcode:   0,
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
