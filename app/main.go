package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

type Message struct {
	id      uint16
	qr      bool
	opcode  uint8
	aa      bool
	tc      bool
	rd      bool
	ra      bool
	z       uint8
	rcode   uint8
	qdcount uint16
	ancount uint16
	nscount uint16
	arcount uint16
}

func (m *Message) bytes() []byte {
	out := []byte{}

	out = binary.BigEndian.AppendUint16(out, m.id)

	var flags uint16
	if m.qr {
		flags |= 1 << 15
	}
	flags |= uint16(m.opcode) << 11
	if m.aa {
		flags |= 1 << 10
	}
	if m.tc {
		flags |= 1 << 9
	}
	if m.rd {
		flags |= 1 << 8
	}
	if m.ra {
		flags |= 1 << 7
	}
	flags |= uint16(m.z) << 4
	flags |= uint16(m.rcode)

	out = binary.BigEndian.AppendUint16(out, flags)
	out = binary.BigEndian.AppendUint16(out, m.qdcount)
	out = binary.BigEndian.AppendUint16(out, m.ancount)
	out = binary.BigEndian.AppendUint16(out, m.nscount)
	out = binary.BigEndian.AppendUint16(out, m.arcount)

	return out
}

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

		headers := Message{
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
		}

		// Create an empty response
		response := headers.bytes()

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
