package lora

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"

	"github.com/Lobaro/slip"
)

type E32PktHdr struct {
	Id   uint32
	Frag uint32
	Last uint32
	Len  byte
}

const (
	HEADERSIZE    = 5
	MAXPACKETSIZE = 58
	PAYLOADSIZE   = MAXPACKETSIZE - HEADERSIZE
)

var Debug bool = false

func NewE32PktHdr(pktid uint32, frag uint32, last uint32, len byte) E32PktHdr {
	var pkt E32PktHdr

	pkt.Id = pktid >> 4
	pkt.Frag = frag
	pkt.Last = last
	pkt.Len = len
	return pkt
}

func (pkt *E32PktHdr) Print() {
	fmt.Printf("Pkt ID: %X\n", pkt.Id)
	fmt.Printf("  frag: %d\n", pkt.Frag)
	fmt.Printf("  last: %d\n", pkt.Last)
	fmt.Printf("   len: %d\n", pkt.Len)
}

func (pkt *E32PktHdr) EncodeHdr() []byte {
	out := make([]byte, 5)
	var idflags uint32 = pkt.Id << 4

	idflags |= (pkt.Frag << 1)
	idflags |= pkt.Last
	a := make([]byte, 4)
	binary.BigEndian.PutUint32(a, idflags)
	out[0] = a[0]
	out[1] = a[1]
	out[2] = a[2]
	out[3] = a[3]
	out[4] = pkt.Len
	//fmt.Printf("Flags: %08b\n", flags)
	return out
}

func DecodeHdr(pkthdr []byte) E32PktHdr {
	var hdr E32PktHdr
	if Debug {
		fmt.Printf("Decode: %X\n", pkthdr)
	}
	hdr.Len = pkthdr[4]
	var y uint32
	_ = binary.Read(bytes.NewReader(pkthdr[0:4]), binary.BigEndian, &y)
	//	fmt.Printf("      %08b\n", pkthdr[3])
	hdr.Id = y >> 4
	hdr.Last = uint32(pkthdr[3] & 1)
	hdr.Frag = 0 | y&2 | (y & 4) | (y & 8)
	hdr.Frag >>= 1
	//	fmt.Printf("Frag: %08b\n", hdr.Frag)
	if Debug {
		hdr.Print()
	}
	//fmt.Printf("hdr: %X\n", uint32(i))
	return hdr
}

func GetSlipHdr(randid, pcount, islast, len uint32) []byte {
	hdrbuf := &bytes.Buffer{}
	pkt := NewE32PktHdr(randid, pcount, islast, byte(len)) // Size is fake
	hdrwriter := slip.NewWriter(hdrbuf)
	hdrwriter.WritePacket(pkt.EncodeHdr())
	return hdrbuf.Bytes()
}
func WriteBytes(p []byte) []byte {
	debugbuf := &bytes.Buffer{}

	rand.Seed(time.Now().UnixNano())
	randid := rand.Uint32()
	buf := &bytes.Buffer{}
	writer := slip.NewWriter(buf)
	err := writer.WritePacket(p)
	if err != nil {
		log.Fatalln(err)
	}
	bufcontent := buf.Bytes()[0 : buf.Len()-1]
	numpackets := len(bufcontent) / (PAYLOADSIZE - 1)

	if Debug {
		println("encoded bytes: ", buf.Len())
		println(numpackets)
	}

	pcount := 0
	fcount := 0
	outbuf := &bytes.Buffer{}
	//islast := 0
	hdrbytes := GetSlipHdr(randid, 0, 0, 0) // Size is fake
	hdrlen := len(hdrbytes) - 1
	payloadsize := MAXPACKETSIZE - hdrlen
	buf.ReadByte() // Skip first slip.ESC
	for {
		if pcount == 0 {
			outbuf = &bytes.Buffer{}
			outbuf.Write(hdrbytes[0:hdrlen]) // write fake header
		}
		pbyte, err := buf.ReadByte()
		if err == io.EOF {
			if Debug {
				println("all read", pcount, outbuf.Len())
			}
			rhdrbytes := GetSlipHdr(randid, uint32(fcount), 1, uint32(pcount-1))
			pkt := outbuf.Bytes()
			pkt[hdrlen-1] = rhdrbytes[hdrlen-1]
			pkt[hdrlen-2] = rhdrbytes[hdrlen-2]
			if Debug {
				fmt.Printf("PKT: %X\n", pkt)
			}
			debugbuf.Write(pkt)
			break
		}
		pcount++
		outbuf.WriteByte(pbyte)
		if pcount == payloadsize-1 {
			islast := 0
			peek, _ := buf.ReadByte()
			if pbyte == slip.END || peek == slip.END {
				islast = 1
				pcount--
				if peek == slip.END {
					outbuf.WriteByte(slip.END)
					pcount++
				}
			} else {
				buf.UnreadByte()
				outbuf.WriteByte(slip.END)
			}

			if Debug {
				println("packet full", pcount, outbuf.Len(), pbyte)
			}

			rhdrbytes := GetSlipHdr(randid, uint32(fcount), uint32(islast), uint32(pcount))
			pkt := outbuf.Bytes()
			pkt[hdrlen-1] = rhdrbytes[hdrlen-1]
			pkt[hdrlen-2] = rhdrbytes[hdrlen-2]
			if Debug {
				fmt.Printf("PKT (%d): %X\n", len(outbuf.Bytes()), outbuf.Bytes())
			}
			debugbuf.Write(pkt)
			pcount = 0
			fcount++
			if islast == 1 {
				break
			}
		}
	}
	return debugbuf.Bytes()
}

// Just for testing...
func ReadBytes(data []byte, inp chan []byte) map[uint32]*bytes.Buffer {
	m := make(map[uint32]*bytes.Buffer)

	debugbuf := &bytes.Buffer{}
	reader := slip.NewReader(bytes.NewReader(data))
	for {
		packet, _, err := reader.ReadPacket()
		if err == io.EOF {
			break
		}
		println("new packet")
		//fmt.Println(packet, isPrefix, err, len(packet))
		hdr := DecodeHdr(packet[0:HEADERSIZE])
		if int(hdr.Len) != len(packet)-HEADERSIZE {
			println(hdr.Len)
			println(len(packet[HEADERSIZE:]))
			log.Fatalln("Packet corrupted! Invalid length!")
			continue
		}
		if m[hdr.Id] == nil {
			m[hdr.Id] = &bytes.Buffer{}
		}
		fmt.Printf("%d %X\n", len(packet[HEADERSIZE:]), packet[HEADERSIZE:])
		debugbuf.Write(packet[HEADERSIZE:])
		m[hdr.Id].Write(packet[HEADERSIZE:])
		if hdr.Last == 1 {
			println("Pkt complete!\n")
		}

	}
	//return debugbuf.Bytes()
	for _, v := range m {
		inp <- v.Bytes()

	}
	return m
}

func RecvLoop(inp chan []byte) {
	port := GetSerial()
	defer port.Close()

	m := make(map[uint32]*bytes.Buffer)
	reader := slip.NewReader(port)
	for {
		packet, _, err := reader.ReadPacket()
		if err != nil {
			println(err)
			continue
		}

		hdr := DecodeHdr(packet[0:HEADERSIZE])
		if int(hdr.Len) != len(packet)-HEADERSIZE {
			fmt.Printf("Bad packet: %X\n", packet)
			println(hdr.Len, len(packet), HEADERSIZE)
			if hdr.Len != 0 {
				println(len(packet[HEADERSIZE:]))
			}
			log.Println("Packet corrupted! Invalid length!")
			continue
		}
		if m[hdr.Id] == nil {
			m[hdr.Id] = &bytes.Buffer{}
		}
		m[hdr.Id].Write(packet[HEADERSIZE:])
		if hdr.Last == 1 {
			inp <- m[hdr.Id].Bytes()
			if Debug {
				fmt.Printf("PKT %X\n", m[hdr.Id].Bytes())
				println("Pkt complete!\n")
			}
			delete(m, hdr.Id)
		}
	}
}
