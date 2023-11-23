package main

import "encoding/binary"

func parse(bytes []byte) *Message {
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

	headers := Headers{
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

	return NewMessage(&headers)
}
