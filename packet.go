package trelay

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Packet interface {
	Length() uint16
	Type() PacketType
	SetType(PacketType) Packet
	Data() []byte

	ResetHead() Packet

	ReadByte() byte
	ReadUint16() uint16
	ReadInt16() int16
	ReadUint32() uint32
	ReadInt32() int32
	ReadUint64() uint64
	ReadInt64() int64
	ReadString() string

	DiscardByte() Packet
	DiscardUint16() Packet
	DiscardInt16() Packet
	DiscardUint32() Packet
	DiscardInt32() Packet
	DiscardUint64() Packet
	DiscardInt64() Packet
	DiscardString() Packet

	PutByte(v byte) Packet
	PutUint16(v uint16) Packet
	PutInt16(v int16) Packet
	PutUint32(v uint32) Packet
	PutInt32(v int32) Packet
	PutUint64(v uint64) Packet
	PutInt64(v int64) Packet
	PutString(v string) Packet
}

type packet struct {
	buf [65535]byte
	ptr uint16
}

func NewPacket(t PacketType) Packet {
	p := &packet{}
	p.ptr = 3
	binary.LittleEndian.PutUint16(p.buf[0:2], 3)
	p.SetType(t)

	return p
}

func NewPacketFromReader(r io.Reader) (Packet, error) {
	var err error
	p := &packet{}
	p.ptr = 3

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
	for p.ptr < p.Length() {
		br, err := r.Read(p.buf[p.ptr:p.Length()])
		if err != nil {
			return p, err
		}

		p.ptr += uint16(br)
	}

	return p, nil
}

func (p *packet) ensurePointerIsAccurate() {
	if p.ptr < 3 {
		p.ptr = 3
	}
}

func (p *packet) ensureCorrectLength() {
	if p.ptr > p.Length() {
		p.addLength(p.ptr - p.Length())
	}
}

func (p *packet) addLength(l uint16) Packet {
	binary.LittleEndian.PutUint16(p.buf[0:2], p.Length()+l)
	return p
}

func (p *packet) Length() uint16 {
	return binary.LittleEndian.Uint16(p.buf[0:2])
}

func (p *packet) Type() PacketType {
	return PacketType(p.buf[2])
}

func (p *packet) SetType(t PacketType) Packet {
	p.buf[2] = byte(t)
	return p
}

func (p *packet) Data() []byte {
	buf := make([]byte, p.Length())
	copy(buf, p.buf[0:p.Length()])
	return buf
}

func (p *packet) ResetHead() Packet {
	p.ptr = 3
	return p
}

// Read functions

func (p *packet) ReadByte() byte {
	p.ensurePointerIsAccurate()
	v := p.buf[p.ptr]
	p.ptr += 1
	return v
}

func (p *packet) ReadUint16() uint16 {
	p.ensurePointerIsAccurate()
	v := binary.LittleEndian.Uint16(p.buf[p.ptr : p.ptr+2])
	p.ptr += 2
	return v
}

func (p *packet) ReadInt16() int16 {
	p.ensurePointerIsAccurate()
	v := int16(binary.LittleEndian.Uint16(p.buf[p.ptr : p.ptr+2]))
	p.ptr += 2
	return v
}

func (p *packet) ReadUint32() uint32 {
	p.ensurePointerIsAccurate()
	v := binary.LittleEndian.Uint32(p.buf[p.ptr : p.ptr+4])
	p.ptr += 4
	return v
}

func (p *packet) ReadInt32() int32 {
	p.ensurePointerIsAccurate()
	v := int32(binary.LittleEndian.Uint32(p.buf[p.ptr : p.ptr+4]))
	p.ptr += 4
	return v
}

func (p *packet) ReadUint64() uint64 {
	p.ensurePointerIsAccurate()
	v := binary.LittleEndian.Uint64(p.buf[p.ptr : p.ptr+8])
	p.ptr += 8
	return v
}

func (p *packet) ReadInt64() int64 {
	p.ensurePointerIsAccurate()
	v := int64(binary.LittleEndian.Uint64(p.buf[p.ptr : p.ptr+8]))
	p.ptr += 8
	return v
}

func (p *packet) ReadString() string {
	p.ensurePointerIsAccurate()
	len := uint16(p.ReadByte())
	if len >= 128 {
		len2 := p.ReadByte()
		len = (len - 128) + uint16(len2<<7) // I have no idea what it does, stolen from popstarfreas/Dimensions
	}

	v := string(p.buf[p.ptr : p.ptr+len])
	p.ptr += len
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

func (p *packet) DiscardInt16() Packet {
	p.ReadInt16()
	return p
}

func (p *packet) DiscardUint32() Packet {
	p.ReadUint32()
	return p
}

func (p *packet) DiscardInt32() Packet {
	p.ReadInt32()
	return p
}

func (p *packet) DiscardUint64() Packet {
	p.ReadUint64()
	return p
}

func (p *packet) DiscardInt64() Packet {
	p.ReadInt64()
	return p
}

func (p *packet) DiscardString() Packet {
	p.ReadString()
	return p
}

// Write functions

func (p *packet) PutByte(v byte) Packet {
	p.ensurePointerIsAccurate()
	p.buf[p.ptr] = v
	p.ptr += 1
	p.ensureCorrectLength()

	return p
}

func (p *packet) PutUint16(v uint16) Packet {
	p.ensurePointerIsAccurate()
	binary.LittleEndian.PutUint16(p.buf[p.ptr:p.ptr+2], v)
	p.ptr += 2
	p.ensureCorrectLength()
	return p
}

func (p *packet) PutInt16(v int16) Packet {
	return p.PutUint16(uint16(v))
}

func (p *packet) PutUint32(v uint32) Packet {
	p.ensurePointerIsAccurate()
	binary.LittleEndian.PutUint32(p.buf[p.ptr:p.ptr+4], v)
	p.ptr += 4
	p.ensureCorrectLength()
	return p
}

func (p *packet) PutInt32(v int32) Packet {
	return p.PutUint32(uint32(v))
}

func (p *packet) PutUint64(v uint64) Packet {
	p.ensurePointerIsAccurate()
	binary.LittleEndian.PutUint64(p.buf[p.ptr:p.ptr+8], v)
	p.ptr += 8
	p.ensureCorrectLength()
	return p
}

func (p *packet) PutInt64(v int64) Packet {
	return p.PutUint64(uint64(v))
}

func (p *packet) PutString(v string) Packet {
	p.ensurePointerIsAccurate()
	strlen := len(v)

	if strlen >= 128 {
		p.PutByte(byte((strlen % 128) + 128))
		p.PutByte(byte(strlen / 128))
	} else {
		p.PutByte(byte(strlen))
	}

	copy(p.buf[p.ptr:], []byte(v))
	p.ptr += uint16(strlen)

	p.ensureCorrectLength()

	return p
}
