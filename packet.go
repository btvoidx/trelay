package trelay

import (
	"encoding/binary"
	"io"
)

type Packet interface {
	// Returns the packet id.
	Id() byte
	// Returns the packet total length.
	Length() uint16
	// Returns a copy of the packet buffer as sent by Terraria (includes packet length and Id).
	Data() []byte
}

// Basic packet implementation
type rawpacket struct{ b []byte }

func (p *rawpacket) Id() byte { return p.b[2] }

func (p *rawpacket) Length() uint16 { return binary.LittleEndian.Uint16(p.b[0:2]) }

func (p *rawpacket) Data() []byte {
	buf := make([]byte, len(p.b))
	copy(buf, p.b)
	return buf
}

// Reads one packet from r. Underlying Packet is *rawpacket.
func ReadPacket(r io.Reader) (Packet, error) {
	ptr := uint16(0)

	// Read packet head (length and id)
	head := make([]byte, 3)
	for ptr < 3 {
		br, err := r.Read(head[ptr:3])
		if err != nil {
			return nil, err
		}
		ptr += uint16(br)
	}

	ln := binary.LittleEndian.Uint16(head[0:2])

	// Read packet data
	b := make([]byte, ln)
	for ptr < ln {
		br, err := r.Read(b[ptr:ln])
		if err != nil {
			return nil, err
		}
		ptr += uint16(br)
	}
	copy(b, head)

	return &rawpacket{b}, nil
}

// TODO
// Parses rawpacket into correct underlying packet struct.
// Returns p if p is not *rawpacket.
// func ParsePacket(p Packet) (Packet, error)
