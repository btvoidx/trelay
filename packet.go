package trelay

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Packet struct {
	ptr uint16
	buf [65535]byte
}

// Reads exactly one packet from io.Reader
func ReadPacket(r io.Reader) (p *Packet, err error) {
	p = &Packet{ptr: 3}

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
	p.ResetHead()

	return p, nil
}

func (p *Packet) ensurePointerIsAccurate() {
	if p.ptr < 3 {
		p.ptr = 3
	}
}

func (p *Packet) Length() uint16 {
	return binary.LittleEndian.Uint16(p.buf[0:2])
}

func (p *Packet) Type() PacketType {
	return PacketType(p.buf[2])
}

func (p *Packet) Data() []byte {
	buf := make([]byte, p.Length())
	copy(buf, p.buf[0:p.Length()])
	return buf
}

func (p *Packet) ResetHead() *Packet {
	p.ptr = 3
	return p
}

// Read functions

func (p *Packet) ReadByte() byte {
	p.ensurePointerIsAccurate()
	v := p.buf[p.ptr]
	p.ptr += 1
	return v
}

func (p *Packet) ReadBytes(l uint16) []byte {
	p.ensurePointerIsAccurate()
	buf := make([]byte, l)
	copy(buf, p.buf[p.ptr:p.ptr+l])
	p.ptr += l
	return buf
}

func (p *Packet) ReadBool() bool {
	v := p.ReadByte()
	return v != 0
}

func (p *Packet) ReadUint16() uint16 {
	p.ensurePointerIsAccurate()
	v := binary.LittleEndian.Uint16(p.buf[p.ptr : p.ptr+2])
	p.ptr += 2
	return v
}

func (p *Packet) ReadInt16() int16 {
	p.ensurePointerIsAccurate()
	v := int16(binary.LittleEndian.Uint16(p.buf[p.ptr : p.ptr+2]))
	p.ptr += 2
	return v
}

func (p *Packet) ReadUint32() uint32 {
	p.ensurePointerIsAccurate()
	v := binary.LittleEndian.Uint32(p.buf[p.ptr : p.ptr+4])
	p.ptr += 4
	return v
}

func (p *Packet) ReadInt32() int32 {
	p.ensurePointerIsAccurate()
	v := int32(binary.LittleEndian.Uint32(p.buf[p.ptr : p.ptr+4]))
	p.ptr += 4
	return v
}

func (p *Packet) ReadUint64() uint64 {
	p.ensurePointerIsAccurate()
	v := binary.LittleEndian.Uint64(p.buf[p.ptr : p.ptr+8])
	p.ptr += 8
	return v
}

func (p *Packet) ReadInt64() int64 {
	p.ensurePointerIsAccurate()
	v := int64(binary.LittleEndian.Uint64(p.buf[p.ptr : p.ptr+8]))
	p.ptr += 8
	return v
}

func (p *Packet) ReadString() string {
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

func (p *Packet) DiscardByte() *Packet {
	p.ReadByte()
	return p
}

func (p *Packet) DiscardBool() *Packet {
	p.DiscardByte()
	return p
}

func (p *Packet) DiscardBytes(l uint16) *Packet {
	p.ReadBytes(l)
	return p
}

func (p *Packet) DiscardUint16() *Packet {
	p.ReadUint16()
	return p
}

func (p *Packet) DiscardInt16() *Packet {
	p.ReadInt16()
	return p
}

func (p *Packet) DiscardUint32() *Packet {
	p.ReadUint32()
	return p
}

func (p *Packet) DiscardInt32() *Packet {
	p.ReadInt32()
	return p
}

func (p *Packet) DiscardUint64() *Packet {
	p.ReadUint64()
	return p
}

func (p *Packet) DiscardInt64() *Packet {
	p.ReadInt64()
	return p
}

func (p *Packet) DiscardString() *Packet {
	p.ReadString()
	return p
}
