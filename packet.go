package trelay

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Packet interface {
	// Packet length is indicated by first 2 bytes, not by buffer length
	Length() uint16
	SetLength(uint16) Packet

	// Packet id is indicated by third byte in buffer
	Id() byte
	SetId(byte) Packet

	// Raw returns a copy of Packet's buffer, not a slice of it
	Raw() []byte

	// Resets pointer to 3.
	ResetHead() Packet

	ReadByte() byte
	ReadUint16() uint16
	ReadString() string

	// Discard functions read from packet, but return packet instead of value.
	// Allows for easy chained skips instead of ignored reads
	DiscardByte() Packet
	DiscardUint16() Packet
	DiscardString() Packet

	PutByte(byte) Packet
	PutString(string) Packet
}

type packet struct {
	buf     [65535]byte
	pointer uint16
}

func NewPacket(id byte) Packet {
	p := &packet{}
	p.pointer = 3
	p.SetId(id)
	p.SetLength(3)

	return p
}

func NewPacketFromReader(r io.Reader) (Packet, error) {
	var err error
	p := &packet{}
	p.pointer = 3

	// Read length
	_, err = r.Read(p.buf[0:2])
	if err != nil {
		return p, err
	}

	if p.Length() < 3 {
		return p, fmt.Errorf("packet is too short: expecting %d bytes, got %d", 3, p.Length())
	}

	if p.Length() > uint16(cap(p.buf)) {
		return p, fmt.Errorf("packet is too long: expecting %d bytes, max size is %d", p.Length(), cap(p.buf))
	}

	// Read ID
	_, err = r.Read(p.buf[2:3])
	if err != nil {
		return p, err
	}

	if p.Length() == 3 {
		return p, err
	}

	// Read data
	for p.pointer < p.Length() {
		br, err := r.Read(p.buf[p.pointer:p.Length()])
		if err != nil {
			return p, err
		}

		p.pointer += uint16(br)
	}

	return p, nil
}

func (p *packet) Length() uint16 {
	return binary.LittleEndian.Uint16(p.buf[0:2])
}

func (p *packet) SetLength(l uint16) Packet {
	binary.LittleEndian.PutUint16(p.buf[0:2], l)
	return p
}

func (p *packet) Id() byte {
	return p.buf[2]
}

func (p *packet) SetId(id byte) Packet {
	p.buf[2] = id
	return p
}

func (p *packet) Raw() []byte {
	buf := make([]byte, p.Length())
	copy(buf, p.buf[0:p.Length()])
	return buf
}

func (p *packet) ResetHead() Packet {
	p.pointer = 3
	return p
}

// Read functions

func (p *packet) ReadByte() byte {
	v := p.buf[p.pointer]
	p.pointer += 1
	return v
}

func (p *packet) ReadUint16() uint16 {
	v := binary.LittleEndian.Uint16(p.buf[p.pointer : p.pointer+2])
	p.pointer += 2
	return v
}

func (p *packet) ReadString() string {
	len := uint16(p.ReadByte())
	if len >= 128 {
		len2 := p.ReadByte()
		len = (len - 128) + uint16(len2<<7) // I have no idea what it does, stolen from popstarfreas/Dimensions
	}

	v := string(p.buf[p.pointer : p.pointer+len])
	p.pointer += len
	return v
}

// Discards

func (p *packet) DiscardByte() Packet {
	p.ReadByte()
	return p
}

func (p *packet) DiscardUint16() Packet {
	p.ReadUint16()
	return p
}

func (p *packet) DiscardString() Packet {
	p.ReadString()
	return p
}

// Write functions

func (p *packet) PutByte(b byte) Packet {
	p.buf[p.pointer] = b
	p.pointer += 1

	p.SetLength(p.Length() + 1)

	return p
}

func (p *packet) PutString(v string) Packet {
	strlen := len(v)

	if strlen >= 128 {
		p.PutByte(byte((strlen % 128) + 128))
		p.PutByte(byte(strlen / 128))
	} else {
		p.PutByte(byte(strlen))
	}

	copy(p.buf[p.pointer:], []byte(v))
	p.pointer += uint16(strlen)

	p.SetLength(p.Length() + uint16(strlen))

	return p
}
