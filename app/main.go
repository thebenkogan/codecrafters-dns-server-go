package main

import (
	"flag"
	"fmt"
)

func main() {
	forwardAddress := flag.String("resolver", "", "address to forward questions")
	flag.Parse()
	fmt.Println("resolver", *forwardAddress)

	client := NewUDPClient("127.0.0.1:2053")
	defer client.Close()

	handler := func(bytes []byte) []byte {
		received := parse(bytes)

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
			response.addAnswer(question.name, 1, 1, 60, [4]uint8{8, 8, 8, 8})
		}

		return response.bytes()
	}

	client.ListenAndServe(handler)
}
