package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

const address string = "127.0.0.1:2053"

func main() {
	forwardAddress := flag.String("resolver", "", "address to forward questions")
	flag.Parse()

	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Panic("Failed to resolve UDP address:", address, err)
	}
	resolverAddr, err := net.ResolveUDPAddr("udp", *forwardAddress)
	if err != nil {
		log.Panic("Failed to resolve UDP address:", *forwardAddress, err)
	}

	udpConn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Panic("Failed to bind to address:", err)
	}
	defer udpConn.Close()

	resolverConn, err := net.DialUDP("udp", nil, resolverAddr)
	if err != nil {
		log.Panic("Failed to dial:", err)
	}
	defer resolverConn.Close()

	buf := make([]byte, 512)
	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			log.Panic("Error receiving data:", err)
		}
		fmt.Printf("Received %d bytes from %s\n", size, source)

		received := parse(buf[:size])
		response := handler(received, resolverConn)

		_, err = udpConn.WriteToUDP(response.bytes(), source)
		if err != nil {
			log.Panic("Failed to send response:", err)
		}
	}
}

func handler(received *Message, resolverConn *net.UDPConn) *Message {
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

	for _, question := range received.questions {
		response.addQuestion(question.name, 1, 1)

		forwardRequest := *received // copies by value
		forwardRequest.questions = []*Question{question}
		forwardRequest.headers.qdcount = 1

		_, err := resolverConn.Write(forwardRequest.bytes())
		if err != nil {
			log.Panic("failed to write to UDP connection", err)
		}

		buf := make([]byte, 512)
		n, err := resolverConn.Read(buf)
		if err != nil {
			log.Panic("failed to read from UDP connection", err)
		}

		forwardResponse := parse(buf[:n])

		if len(forwardResponse.answers) > 0 {
			answer := forwardResponse.answers[0]
			response.addAnswer(answer.name, answer.typ, answer.class, answer.ttl, answer.rdata)
		}
	}

	return response
}
