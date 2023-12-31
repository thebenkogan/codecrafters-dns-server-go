package main

import (
	"encoding/binary"
	"strings"
)

type Message struct {
	headers   *Headers
	questions []*Question
	answers   []*ResourceRecord
}

func NewMessage(headers *Headers) *Message {
	return &Message{
		headers:   headers,
		questions: []*Question{},
		answers:   []*ResourceRecord{},
	}
}

func (m *Message) bytes() []byte {
	out := m.headers.bytes()
	for _, question := range m.questions {
		out = append(out, question.bytes()...)
	}
	for _, answer := range m.answers {
		out = append(out, answer.bytes()...)
	}
	return out
}

func (m *Message) addQuestion(name string, typ uint16, class uint16) {
	m.questions = append(m.questions, &Question{name, typ, class})
	m.headers.qdcount++
}

func (m *Message) addAnswer(rr *ResourceRecord) {
	m.answers = append(m.answers, rr)
	m.headers.ancount++
}

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

func encodeName(name string) []byte {
	out := []byte{}
	for _, label := range strings.Split(name, ".") {
		out = append(out, byte(len(label)))
		out = append(out, []byte(label)...)
	}
	out = append(out, []byte("\x00")...)
	return out
}

type Question struct {
	name  string
	typ   uint16
	class uint16
}

func (q *Question) bytes() []byte {
	out := []byte{}

	out = append(out, encodeName(q.name)...)
	out = binary.BigEndian.AppendUint16(out, q.typ)
	out = binary.BigEndian.AppendUint16(out, q.class)

	return out
}

type ResourceRecord struct {
	name     string
	typ      uint16
	class    uint16
	ttl      uint32
	rdlength uint16
	rdata    [4]uint8
}

func (r *ResourceRecord) bytes() []byte {
	out := []byte{}

	out = append(out, encodeName(r.name)...)
	out = binary.BigEndian.AppendUint16(out, r.typ)
	out = binary.BigEndian.AppendUint16(out, r.class)
	out = binary.BigEndian.AppendUint32(out, r.ttl)
	out = binary.BigEndian.AppendUint16(out, r.rdlength)
	out = append(out, r.rdata[:]...)

	return out
}
