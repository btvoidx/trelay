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

func (w *Builder) Id() byte { return w.buf[2] }

func (w *Builder) Length() uint16 { return uint16(len(w.buf)) }

func (w *Builder) Data() []byte {
	w.setupBuffer()

	buf := make([]byte, w.Length())
	copy(buf, w.buf)
	binary.LittleEndian.PutUint16(buf[0:2], w.Length())

	return buf
}

func (w *Builder) setupBuffer() {
	if len(w.buf) < 3 {
		// Most packets are rather short, so allocating 16 bytes is more than enough at first.
		w.buf = make([]byte, 3, 16)
	}
}

func (w *Builder) SetId(id byte) *Builder {
	w.setupBuffer()
	w.buf[2] = byte(id)
	return w
}

func (w *Builder) WriteBytes(v []byte) *Builder {
	w.setupBuffer()
	w.buf = append(w.buf, v...)
	return w
}

func (w *Builder) WriteByte(v byte) *Builder {
	w.setupBuffer()
	w.buf = append(w.buf, v)
	return w
}

func (w *Builder) WriteBool(v bool) *Builder {
	if v {
		return w.WriteByte(1)
	}
	return w.WriteByte(0)
}

func (w *Builder) WriteUint16(v uint16) *Builder {
	w.setupBuffer()
	vbuf := make([]byte, 2)
	binary.LittleEndian.PutUint16(vbuf, v)
	w.buf = append(w.buf, vbuf...)
	return w
}

func (w *Builder) WriteInt16(v int16) *Builder {
	return w.WriteUint16(uint16(v))
}

func (w *Builder) WriteUint32(v uint32) *Builder {
	w.setupBuffer()
	vbuf := make([]byte, 4)
	binary.LittleEndian.PutUint32(vbuf, v)
	w.buf = append(w.buf, vbuf...)
	return w
}

func (w *Builder) WriteInt32(v int32) *Builder {
	return w.WriteUint32(uint32(v))
}

func (w *Builder) WriteUint64(v uint64) *Builder {
	w.setupBuffer()
	vbuf := make([]byte, 8)
	binary.LittleEndian.PutUint64(vbuf, v)
	w.buf = append(w.buf, vbuf...)
	return w
}

func (w *Builder) WriteInt64(v int64) *Builder {
	return w.WriteUint64(uint64(v))
}

func (w *Builder) WriteFloat32(v float32) *Builder {
	return w.WriteUint32(math.Float32bits(v))
}

func (w *Builder) WriteFloat64(v float64) *Builder {
	return w.WriteUint64(math.Float64bits(v))
}

func (w *Builder) WriteString(v string) *Builder {
	w.setupBuffer()
	if l := len(v); l >= 128 {
		w.WriteByte(byte((l % 128) + 128))
		w.WriteByte(byte(l / 128))
	} else {
		w.WriteByte(byte(l))
	}
	w.buf = append(w.buf, v...)
	return w
}

// A utility function around w.WriteString(fmt.Sprintf(format, args...))
func (w *Builder) WriteStringf(format string, a ...any) *Builder {
	return w.WriteString(fmt.Sprintf(format, a...))
}
