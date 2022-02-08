package trelay

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Packet struct {
	ptr uint16
	buf []byte
}

// Reads exactly one packet from io.Reader. Packet pointer is nil if an error occurs even if most of the packet was read successfully
func ReadPacket(r io.Reader) (*Packet, error) {
	p := &Packet{ptr: 3, buf: make([]byte, 3)}

	// Read length
	_, err := r.Read(p.buf[0:2])
	if err != nil {
		return p, err
	}

	if p.Length() < 3 || p.Length() > 65535 {
		return nil, fmt.Errorf("bad packet length")
	}

	// Read type
	_, err = r.Read(p.buf[2:3])
	if err != nil {
		return p, err
	}

	if p.Length() > 3 {
		buf := make([]byte, p.Length())
		copy(buf[0:3], p.buf[0:3])
		p.buf = buf
	}

	// Read data
	for p.ptr < p.Length() {
		br, err := r.Read(p.buf[p.ptr:p.Length()])
		if err != nil {
			return nil, err
		}

		p.ptr += uint16(br)
	}

	p.ResetHead()

	return p, nil
}

func (p Packet) String() string {
	return fmt.Sprintf("{type:%d, len:%d}", p.Type(), p.Length())
}

func (p *Packet) ensureAccuratePointer() {
	if p.ptr < 3 {
		p.ptr = 3
	}
}

func (p *Packet) canReadN(l uint16) bool {
	return p.ptr+l <= p.Length()
}

func (p *Packet) Length() uint16 {
	return binary.LittleEndian.Uint16(p.buf[0:2])
}

func (p *Packet) Type() PacketType {
	return PacketType(p.buf[2])
}

// Clones packets internal buffer and returns it
func (p *Packet) Data() []byte {
	buf := make([]byte, p.Length())
	copy(buf, p.buf[0:p.Length()])
	return buf
}

// Resets head to 3 (start of packet body)
func (p *Packet) ResetHead() {
	p.ptr = 3
}

// Advances head l bytes, returns io.EOF if unsuccessful
func (p *Packet) AdvanceHead(l uint16) error {
	p.ensureAccuratePointer()
	if !p.canReadN(l) {
		return io.EOF
	}
	p.ptr += l
	return nil
}

// Read functions

// Reads and returns a byte, error is io.EOF if unsuccessful
func (p *Packet) ReadByte() (byte, error) {
	p.ensureAccuratePointer()
	if !p.canReadN(1) {
		return 0, io.EOF
	}
	v := p.buf[p.ptr]
	p.ptr += 1
	return v, nil
}

// Reads and returns a byte, panics if unsuccessful
func (p *Packet) MustReadByte() byte {
	v, err := p.ReadByte()
	if err != nil {
		panic(err)
	}
	return v
}

// Reads a byte and returns true if is is != 0, error is io.EOF if unsuccessful
func (p *Packet) ReadBool() (bool, error) {
	v, err := p.ReadByte()
	if err != nil {
		return false, err
	}
	return v != 0, nil
}

// Reads a byte and returns true if is is != 0, panics if unsuccessful
func (p *Packet) MustReadBool() bool {
	v, err := p.ReadBool()
	if err != nil {
		panic(err)
	}
	return v
}

// Reads and returns l bytes, error is io.EOF if unsuccessful. Head is not advanced on error, so it is still possible to read a smaller value from packet
func (p *Packet) ReadBytes(l uint16) ([]byte, error) {
	p.ensureAccuratePointer()
	if !p.canReadN(l) {
		return nil, io.EOF
	}
	buf := make([]byte, l)
	copy(buf, p.buf[p.ptr:p.ptr+l])
	p.ptr += l
	return buf, nil
}

// Reads and returns l bytes, panics if unsuccessful
func (p *Packet) MustReadBytes(l uint16) []byte {
	v, err := p.ReadBytes(l)
	if err != nil {
		panic(err)
	}
	return v
}

// Reads and returns uint16, error is io.EOF if unsuccessful. Head is not advanced on error, so it is still possible to read a smaller value from packet
func (p *Packet) ReadUint16() (uint16, error) {
	p.ensureAccuratePointer()
	if !p.canReadN(2) {
		return 0, io.EOF
	}
	v := binary.LittleEndian.Uint16(p.buf[p.ptr : p.ptr+2])
	p.ptr += 2
	return v, nil
}

// Reads and returns uint16, panics if unsuccessful
func (p *Packet) MustReadUint16() uint16 {
	v, err := p.ReadUint16()
	if err != nil {
		panic(err)
	}
	return v
}

// Reads and returns int16, error is io.EOF if unsuccessful. Head is not advanced on error, so it is still possible to read a smaller value from packet
func (p *Packet) ReadInt16() (int16, error) {
	v, err := p.ReadUint16()
	return int16(v), err
}

// Reads and returns int16, panics if unsuccessful
func (p *Packet) MustReadInt16() int16 {
	v, err := p.ReadInt16()
	if err != nil {
		panic(err)
	}
	return v
}

// Reads and returns uint32, error is io.EOF if unsuccessful. Head is not advanced on error, so it is still possible to read a smaller value from packet
func (p *Packet) ReadUint32() (uint32, error) {
	p.ensureAccuratePointer()
	if !p.canReadN(4) {
		return 0, io.EOF
	}
	v := binary.LittleEndian.Uint32(p.buf[p.ptr : p.ptr+4])
	p.ptr += 4
	return v, nil
}

// Reads and returns uint32, panics if unsuccessful
func (p *Packet) MustReadUint32() uint32 {
	v, err := p.ReadUint32()
	if err != nil {
		panic(err)
	}
	return v
}

// Reads and returns int32, error is io.EOF if unsuccessful. Head is not advanced on error, so it is still possible to read a smaller value from packet
func (p *Packet) ReadInt32() (int32, error) {
	v, err := p.ReadUint16()
	return int32(v), err
}

// Reads and returns int32, panics if unsuccessful
func (p *Packet) MustReadInt32() int32 {
	v, err := p.ReadInt32()
	if err != nil {
		panic(err)
	}
	return v
}

// Reads and returns uint64, error is io.EOF if unsuccessful. Head is not advanced on error, so it is still possible to read a smaller value from packet
func (p *Packet) ReadUint64() (uint64, error) {
	p.ensureAccuratePointer()
	if !p.canReadN(8) {
		return 0, io.EOF
	}
	v := binary.LittleEndian.Uint64(p.buf[p.ptr : p.ptr+8])
	p.ptr += 8
	return v, nil
}

// Reads and returns uint64, panics if unsuccessful
func (p *Packet) MustReadUint64() uint64 {
	v, err := p.ReadUint64()
	if err != nil {
		panic(err)
	}
	return v
}

// Reads and returns int64, error is io.EOF if unsuccessful. Head is not advanced on error, so it is still possible to read a smaller value from packet
func (p *Packet) ReadInt64() (int64, error) {
	v, err := p.ReadUint16()
	return int64(v), err
}

// Reads and returns int64, panics if unsuccessful
func (p *Packet) MustReadInt64() int64 {
	v, err := p.ReadInt64()
	if err != nil {
		panic(err)
	}
	return v
}

// Reads and returns a string, error is io.EOF if unsuccessful. Head is not advanced on error, so it is still possible to read a smaller value from packet
func (p *Packet) ReadString() (string, error) {
	p.ensureAccuratePointer()
	// saves current head so it can be restored on error
	prevptr := p.ptr

	len1, err := p.ReadByte()
	if err != nil {
		return "", err
	}
	len := uint16(len1)

	if len1 >= 128 {
		len2, err := p.ReadByte()
		if err != nil {
			p.ptr = prevptr
			return "", err
		}
		len = (len - 128) + uint16(len2<<7) // I have no idea what it does, stolen from popstarfreas/Dimensions
	}

	if !p.canReadN(len) {
		p.ptr = prevptr
		return "", io.EOF
	}

	v := string(p.buf[p.ptr : p.ptr+len])
	p.ptr += len
	return v, nil
}

// Reads and returns a string, panics if unsuccessful
func (p *Packet) MustReadString() string {
	v, err := p.ReadString()
	if err != nil {
		panic(err)
	}
	return v
}
