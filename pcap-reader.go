package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"unsafe"
)

type PCapPacketHeader struct {
	TimestampSeconds   uint32
	TimestampExtra     uint32
	CapturedDataLength uint32
	TotalDataLength    uint32
}

type PCapPacket struct {
	Header  PCapPacketHeader
	RawData []byte
}

func (p PCapPacket) ToString() string {
	extraUnits := "ns"
	return fmt.Sprintf("%x s, %x %s, %x, %x", p.Header.TimestampSeconds, p.Header.TimestampExtra, extraUnits, p.Header.CapturedDataLength, p.Header.TotalDataLength)
}

type PCapSavefileHeader struct {
	MagicNumber         uint32
	MajorVersion        uint16
	MinorVersion        uint16
	TZOffset            uint32
	TZAccuracy          uint32
	SnapshotLength      uint32
	LinkLayerHeaderType uint32
}

func (h PCapSavefileHeader) ToString() string {
	return fmt.Sprintf("%x, %x, %x, %x, %x, %x, %x", h.MagicNumber, h.MajorVersion, h.MinorVersion, h.TZOffset, h.TZAccuracy, h.SnapshotLength, h.LinkLayerHeaderType)
}

type PCapSavefile struct {
	Header  PCapSavefileHeader
	Packets []*PCapPacket
}

type PCapReader struct {
	position  int
	byteOrder binary.ByteOrder //
	bytes     []byte
}

func (r PCapReader) RemainingBytes() int {
	return len(r.bytes) - r.position
}

func (r *PCapReader) ReadPCapGlobalHeader() PCapSavefileHeader {
	var result PCapSavefileHeader
	buf := bytes.NewReader(r.bytes[0:24])
	err := binary.Read(buf, r.byteOrder, &result)
	if err != nil {
		panic(err)
	}
	r.position += 24
	return result
}

func (r *PCapReader) ReadPCapPacket() *PCapPacket {
	var header PCapPacketHeader
	buf := bytes.NewReader(r.bytes[r.position : r.position+16])
	err := binary.Read(buf, r.byteOrder, &header)
	if err != nil {
		panic(err)
	}
	r.position += 16
	data := r.bytes[r.position : r.position+int(header.CapturedDataLength)]
	r.position += int(header.CapturedDataLength)
	return &PCapPacket{header, data}
}

func MakePCapReader(contents []byte) *PCapReader {
	fileByteOrder := hostByteOrder()
	if !fileByteOrderMatchesHost(contents) {
		if fileByteOrder.String() == binary.BigEndian.String() {
			fileByteOrder = binary.LittleEndian
		} else {
			fileByteOrder = binary.BigEndian
		}
	}
	return &PCapReader{0, fileByteOrder, contents}
}

func fileByteOrderMatchesHost(contents []byte) bool {
	var magicNumber uint32
	byteOrder := hostByteOrder()
	buf := bytes.NewReader(contents[0:4])
	err := binary.Read(buf, byteOrder, &magicNumber)
	if err != nil {
		panic(err)
	}
	switch magicNumber {
	case 0xa1b2c3d4, 0xa1b23c4d:
		return true
	case 0xd4c3b2a1, 0x4d3cb2a1:
		return false
	default:
		panic("Unexpected magic number")
	}
}

func hostByteOrder() binary.ByteOrder {
	buf := [2]byte{}
	*(*uint32)(unsafe.Pointer(&buf[0])) = uint32(0xABCD)
	switch buf {
	case [2]byte{0xCD, 0xAB}:
		return binary.LittleEndian
	case [2]byte{0xAB, 0xCD}:
		return binary.BigEndian
	default:
		panic("Could not determine native endianness.")
	}
}
