package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type MACAddress [6]byte

func (m MACAddress) String() string {
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", m[0], m[1], m[2], m[3], m[4], m[5])
}

type DataLinkEthernetHeader struct {
	Destination MACAddress
	Source      MACAddress
	EtherType   uint16
}

type EthernetFrame struct {
	Header  DataLinkEthernetHeader
	Payload []byte
}

func ReadEthernetFrame(data []byte) *EthernetFrame {
	header := ParseEthernetHeader(data)
	return &EthernetFrame{header, data[14:]}
}

func ParseEthernetHeader(data []byte) DataLinkEthernetHeader {
	var header DataLinkEthernetHeader
	buf := bytes.NewReader(data)
	err := binary.Read(buf, binary.BigEndian, &header)
	if err != nil {
		panic(err)
	}
	return header
}
