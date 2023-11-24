package main

import (
	"fmt"
	"log"
	"net"
)

type Client struct {
	conn *net.UDPConn
}

func NewUDPClient(address string) *Client {
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Panic("Failed to resolve UDP address:", address, err)
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Panic("Failed to bind to address:", err)
	}

	return &Client{udpConn}
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) ListenAndServe(handler func(bytes []byte) []byte) {
	buf := make([]byte, 512)

	for {
		size, source, err := c.conn.ReadFromUDP(buf)
		if err != nil {
			log.Panic("Error receiving data:", err)
		}
		fmt.Printf("Received %d bytes from %s\n", size, source)

		_, err = c.conn.WriteToUDP(handler(buf[:size]), source)
		if err != nil {
			log.Panic("Failed to send response:", err)
		}
	}
}

func (c *Client) SendAndReceive(packet []byte) []byte {
	_, err := c.conn.Write(packet)
	if err != nil {
		log.Panic("Failed to send packet:", err)
	}

	buf := make([]byte, 512)
	size, _, err := c.conn.ReadFromUDP(buf)
	if err != nil {
		log.Panic("Error receiving data:", err)
	}

	return buf[:size]
}
