package main

import ()

type DataLinkEthernetHeader struct {
	Destination uint64
	Source      uint64
	EtherType   uint16
}

type EthernetFrame struct {
	Header             DataLinkEthernetHeader
	Payload            []byte
	FrameCheckSequence uint32
}

func ParseEthernetHeader(data []byte) DataLinkEthernetHeader {

}
