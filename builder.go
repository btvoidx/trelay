package trelay

import (
	"encoding/binary"
	"fmt"
	"math"
)

// Builder is used to efficiently build a packet.
// Do not copy a non-zero Builder, as buffer will be shared.
// Zero value is ready to use.
// Builder implemets Packet inteface.
type Builder struct {
	buf []byte
}

var _ Packet = (*Builder)(nil)

func (b *Builder) Id() byte { return b.buf[2] }

func (b *Builder) Length() uint16 { return uint16(len(b.buf)) }

func (b *Builder) Data() []byte {
	b.setupBuffer()

	buf := make([]byte, b.Length())
	copy(buf, b.buf)
	binary.LittleEndian.PutUint16(buf[0:2], b.Length())

	return buf
}

func (b *Builder) setupBuffer() {
	if len(b.buf) < 3 {
		// Most packets are rather short, so allocating 16 bytes is more than enough at first.
		b.buf = make([]byte, 3, 16)
	}
}

func (b *Builder) SetId(id byte) *Builder {
	b.setupBuffer()
	b.buf[2] = byte(id)
	return b
}

func (b *Builder) WriteBytes(v []byte) *Builder {
	b.setupBuffer()
	b.buf = append(b.buf, v...)
	return b
}

func (b *Builder) WriteByte(v byte) *Builder {
	b.setupBuffer()
	b.buf = append(b.buf, v)
	return b
}

func (b *Builder) WriteBool(v bool) *Builder {
	if v {
		return b.WriteByte(1)
	}
	return b.WriteByte(0)
}

func (b *Builder) WriteUint16(v uint16) *Builder {
	b.setupBuffer()
	vbuf := make([]byte, 2)
	binary.LittleEndian.PutUint16(vbuf, v)
	b.buf = append(b.buf, vbuf...)
	return b
}

func (b *Builder) WriteInt16(v int16) *Builder {
	return b.WriteUint16(uint16(v))
}

func (b *Builder) WriteUint32(v uint32) *Builder {
	b.setupBuffer()
	vbuf := make([]byte, 4)
	binary.LittleEndian.PutUint32(vbuf, v)
	b.buf = append(b.buf, vbuf...)
	return b
}

func (b *Builder) WriteInt32(v int32) *Builder {
	return b.WriteUint32(uint32(v))
}

func (b *Builder) WriteUint64(v uint64) *Builder {
	b.setupBuffer()
	vbuf := make([]byte, 8)
	binary.LittleEndian.PutUint64(vbuf, v)
	b.buf = append(b.buf, vbuf...)
	return b
}

func (b *Builder) WriteInt64(v int64) *Builder {
	return b.WriteUint64(uint64(v))
}

func (b *Builder) WriteFloat32(v float32) *Builder {
	return b.WriteUint32(math.Float32bits(v))
}

func (b *Builder) WriteFloat64(v float64) *Builder {
	return b.WriteUint64(math.Float64bits(v))
}

func (b *Builder) WriteString(v string) *Builder {
	b.setupBuffer()
	if l := len(v); l >= 128 {
		b.WriteByte(byte((l % 128) + 128))
		b.WriteByte(byte(l / 128))
	} else {
		b.WriteByte(byte(l))
	}
	b.buf = append(b.buf, v...)
	return b
}

// A utility function around w.WriteString(fmt.Sprintf(format, args...))
func (b *Builder) WriteStringf(format string, a ...any) *Builder {
	return b.WriteString(fmt.Sprintf(format, a...))
}
