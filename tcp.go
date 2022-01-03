package main

type TCPOption struct {
	Kind   uint8
	Length uint8
	Data   []byte
}

type TCPHeader struct {
	SourcePort            uint16
	DestinationPort       uint16
	SequenceNumber        uint32
	AcknowledgementNumber uint32
	DataOffset            uint8 // 4 bits
	NS                    bool  //
	CWR                   bool
	ECE                   bool
	URG                   bool
	ACK                   bool
	PSH                   bool
	RST                   bool
	SYN                   bool
	FIN                   bool
	WindowSize            uint16
	Checksum              uint16
	UrgentPointer         uint16
	Options               []TCPOption
}

type TCPSegment struct {
	Header TCPHeader
	Data   []byte
}

func ReadTCPSegment(data []byte) *TCPSegment {
	header := ParseTCPHeader(data)
	headerSizeBytes := 4 * int(header.DataOffset)
	return &TCPSegment{header, data[headerSizeBytes:]}
}

func ParseTCPHeader(data []byte) TCPHeader {
	h := TCPHeader{}
	h.SourcePort = (uint16(data[0]) << 8) | uint16(data[1])
	h.DestinationPort = (uint16(data[2]) << 8) | uint16(data[3])
	h.SequenceNumber = uint32(data[4])<<24 | uint32(data[5])<<16 | uint32(data[6])<<8 | uint32(data[7])
	h.AcknowledgementNumber = uint32(data[8])<<24 | uint32(data[9])<<16 | uint32(data[10])<<8 | uint32(data[11])
	h.DataOffset = data[12] >> 4
	if data[12]&0x1 == 1 {
		h.NS = true
	}
	if (data[13]>>7)&1 == 1 {
		h.CWR = true
	}
	if (data[13]>>6)&1 == 1 {
		h.ECE = true
	}
	if (data[13]>>5)&1 == 1 {
		h.URG = true
	}
	if (data[13]>>4)&1 == 1 {
		h.ACK = true
	}
	if (data[13]>>3)&1 == 1 {
		h.PSH = true
	}
	if (data[13]>>2)&1 == 1 {
		h.RST = true
	}
	if (data[13]>>1)&1 == 1 {
		h.SYN = true
	}
	if data[13]&1 == 1 {
		h.FIN = true
	}
	h.WindowSize = uint16(data[14])<<8 | uint16(data[15])
	h.Checksum = uint16(data[16])<<8 | uint16(data[17])
	h.UrgentPointer = uint16(data[19])<<8 | uint16(data[20])
	h.Options = ReadOptions(data, h.DataOffset)
	return h
}

func ReadOptions(data []byte, offset uint8) []TCPOption {
	options := []TCPOption{}
	end := int(4 * (offset - 5))
	for pos := 21; pos < end; {
		kind := data[pos]
		switch kind {
		case 0:
			return options
		case 1:
			pos++
		case 2, 3, 4, 5, 8:
			length := data[pos+1]
			d := make([]byte, length-2)
			for j := 2; j < int(length); j++ {
				d[j-2] = data[pos+j]
			}
			o := TCPOption{Kind: 2, Length: length, Data: d}
			options = append(options, o)
			pos += int(length)
		default:
			panic("Unknown TCP option!")
		}
	}
	return options
}
