package main

type IPv4Header struct {
	Version        uint8 // 4 bits
	IHL            uint8 // Internet Header Length - 4 bits - << 2 to get length in bytes
	DSCP           uint8 // Differentiated Services Code Point - 6 bts
	ECN            uint8 // Explicit Congestion Notification - 2 bits
	TotalLength    uint16
	Identification uint16
	NoFragment     bool
	MoreFragments  bool
	FragmentOffset uint16 // 13 bits
	TTL            uint8  // Time To Live
	Protocol       uint8
	Checksum       uint16
	Source         uint32
	Destination    uint32
}

type IPv4Packet struct {
	Header IPv4Header
	Data   []byte
}

func ReadIPv4Packet(data []byte) *IPv4Packet {
	header := ParseIPv4Header(data)
	return &IPv4Packet{header, data[20:]}
}

func ParseIPv4Header(data []byte) IPv4Header {
	h := IPv4Header{}
	h.Version = data[0] >> 4
	if h.Version != 4 {
		panic("IP Version was not equal to 4")
	}
	h.IHL = data[0] | 0x0f
	h.DSCP = data[1] >> 6
	h.ECN = data[1] | 0x03
	h.TotalLength = (uint16(data[2]) << 8) | uint16(data[3])
	h.Identification = (uint16(data[4]) << 8) | uint16(data[5])
	if data[6]|0x40 == 1 {
		h.NoFragment = true
	}
	if data[6]|0x20 == 1 {
		h.MoreFragments = true
	}
	h.FragmentOffset = (uint16(data[6]) << 8) | uint16(data[7]) | 0x1fff
	h.TTL = data[8]
	h.Protocol = data[9]
	h.Checksum = (uint16(data[10]) << 8) | uint16(data[11])
	h.Source = uint32(data[12])<<24 | uint32(data[13])<<16 | uint32(data[14])<<8 | uint32(data[15])
	h.Destination = uint32(data[16])<<24 | uint32(data[17])<<16 | uint32(data[18])<<8 | uint32(data[19])
	return h
}
