package main

import "encoding/binary"

type Headers struct {
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

func (h *Headers) flags() uint16 {
	var flags uint16
	if h.qr {
		flags |= 1 << 15
	}
	flags |= uint16(h.opcode) << 11
	if h.aa {
		flags |= 1 << 10
	}
	if h.tc {
		flags |= 1 << 9
	}
	if h.rd {
		flags |= 1 << 8
	}
	if h.ra {
		flags |= 1 << 7
	}
	flags |= uint16(h.z) << 4
	flags |= uint16(h.rcode)

	return flags
}

func (h *Headers) bytes() []byte {
	out := []byte{}
	out = binary.BigEndian.AppendUint16(out, h.id)
	out = binary.BigEndian.AppendUint16(out, h.flags())
	out = binary.BigEndian.AppendUint16(out, h.qdcount)
	out = binary.BigEndian.AppendUint16(out, h.ancount)
	out = binary.BigEndian.AppendUint16(out, h.nscount)
	out = binary.BigEndian.AppendUint16(out, h.arcount)
	return out
}
