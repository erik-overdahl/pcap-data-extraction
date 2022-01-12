package main

import (
	"bytes"
	"fmt"
	"os"
	"sort"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("FILE is a required argument")
		os.Exit(0)
	} else if len(os.Args) > 2 {
		fmt.Println("Too many arguments; 1 expected")
		os.Exit(1)
	}
	filepath := os.Args[1]
	contents := readFile(filepath)
	file := ReadPCapFile(contents)
	ethernetFrames := make([]*EthernetFrame, len(file.Packets))
	for i, p := range file.Packets {
		ethernetFrames[i] = ReadEthernetFrame(p.RawData)
	}

	fmt.Printf("Found %d packets\n", len(file.Packets))
	PrintEthernetHeadersInfo(ethernetFrames)

	ipPackets := make([]*IPv4Packet, len(file.Packets))
	for i, f := range ethernetFrames {
		ipPackets[i] = ReadIPv4Packet(f.Payload)
	}

	// the client has to initiate the connection
	// this assumes that the capture file is complete and
	// that the first packet is the start of the connection
	clientIp := ipPackets[0].Header.Source

	serverResponses := []*TCPSegment{}
	seqNums := make(map[uint32]bool)
	for _, p := range ipPackets {
		s := ReadTCPSegment(p.Data)
		if !seqNums[s.Header.SequenceNumber] && p.Header.Destination == clientIp {
			seqNums[s.Header.SequenceNumber] = true
			serverResponses = append(serverResponses, s)
		}
	}

	sort.Slice(serverResponses, func(i, j int) bool {
		return serverResponses[i].Header.SequenceNumber < serverResponses[j].Header.SequenceNumber
	})

	values := make([][]byte, len(serverResponses))
	for i, r := range serverResponses {
		values[i] = r.Data
	}
	rawHttpResponse := bytes.Join(values, []byte{})
	pieces := bytes.SplitN(rawHttpResponse, []byte{'\r', '\n', '\r', '\n'}, 2)
	fmt.Printf("%s\n", string(pieces[0]))
	err := os.WriteFile("out.jpeg", pieces[1], 0644)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Wrote image to file!\n")
}

func PrintEthernetHeadersInfo(frames []*EthernetFrame) {
	ipVersion := make(map[uint16]int)
	destinations := make(map[MACAddress]int)
	sources := make(map[MACAddress]int)
	for _, frame := range frames {
		if _, exists := ipVersion[frame.Header.EtherType]; !exists {
			ipVersion[frame.Header.EtherType] = 1
		} else {
			ipVersion[frame.Header.EtherType] += 1
		}
		if _, exists := destinations[frame.Header.Destination]; !exists {
			destinations[frame.Header.Destination] = 1
		} else {
			destinations[frame.Header.Destination] += 1
		}
		if _, exists := sources[frame.Header.Source]; !exists {
			sources[frame.Header.Source] = 1
		} else {
			sources[frame.Header.Source] += 1
		}
	}
	for v, count := range ipVersion {
		fmt.Printf("%d have EtherType %#x\n", count, v)
	}
	for v, count := range destinations {
		fmt.Printf("%d have destination %s\n", count, v.String())
	}
	for v, count := range sources {
		fmt.Printf("%d have source %s\n", count, v.String())
	}
}

func ReadPCapFile(contents []byte) *PCapSavefile {
	r := MakePCapReader(contents)
	header := r.ReadPCapGlobalHeader()
	file := &PCapSavefile{header, []*PCapPacket{}}
	for r.RemainingBytes() > 16 {
		packet := r.ReadPCapPacket()
		file.Packets = append(file.Packets, packet)
	}
	return file
}

func readFile(filepath string) []byte {
	contents, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	return contents
}
