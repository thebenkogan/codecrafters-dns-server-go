package main

import (
	"encoding/binary"
	"strings"
)

func parse(bytes []byte) *Message {
	headers := parseHeaders(bytes[:12])
	questions, _ := parseQuestions(bytes[12:], int(headers.qdcount))
	message := NewMessage(headers)
	message.questions = questions
	return message
}

func parseHeaders(bytes []byte) *Headers {
	id := binary.BigEndian.Uint16(bytes[:2])
	flags := binary.BigEndian.Uint16(bytes[2:4])
	qdcount := binary.BigEndian.Uint16(bytes[4:6])
	ancount := binary.BigEndian.Uint16(bytes[6:8])
	nscount := binary.BigEndian.Uint16(bytes[8:10])
	arcount := binary.BigEndian.Uint16(bytes[10:12])

	qr := ((1 << 15) & flags) != 0
	opcode := uint8((flags >> 11) & 0b1111)
	aa := ((1 << 10) & flags) != 0
	tc := ((1 << 9) & flags) != 0
	rd := ((1 << 8) & flags) != 0
	ra := ((1 << 7) & flags) != 0
	z := uint8((flags >> 4) & 0b111)
	rcode := uint8(flags & 0b1111)

	return &Headers{
		id,
		qr,
		opcode,
		aa,
		tc,
		rd,
		ra,
		z,
		rcode,
		qdcount,
		ancount,
		nscount,
		arcount,
	}
}

func parseName(bytes []byte) (string, []byte) {
	labels := []string{}

	currIndex := 0
	for bytes[currIndex] != 0 {
		length := int(bytes[currIndex])
		label := string(bytes[currIndex+1 : currIndex+length+1])
		labels = append(labels, label)
		currIndex = currIndex + length + 1
	}

	return strings.Join(labels, "."), bytes[currIndex+1:]
}

func parseQuestions(bytes []byte, amount int) ([]Question, []byte) {
	questions := make([]Question, amount)

	remainingBytes := bytes
	for i := 0; i < amount; i++ {
		name, remaining := parseName(remainingBytes)
		typ := binary.BigEndian.Uint16(remaining[0:2])
		class := binary.BigEndian.Uint16(remaining[2:4])
		remainingBytes = remaining[4:]
		questions[i] = Question{name, typ, class}
	}

	return questions, remainingBytes
}
