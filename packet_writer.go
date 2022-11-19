package trelay

import (
	"encoding/binary"
	"fmt"
	"math"
)

// Writer is used to efficiently build a terraria tcp packet.
// Do not copy a non-zero Writer, as buffer will be shared.
type Writer struct {
	buf []byte
}

// Creates a copy of the buffer and returns it.
func (w *Writer) Data() []byte {
	w.setupBuffer()
	l := len(w.buf)
	buf := make([]byte, l)
	copy(buf, w.buf)
	binary.LittleEndian.PutUint16(buf[0:2], uint16(l))
	return buf
}

func (w *Writer) setupBuffer() {
	if len(w.buf) < 3 {
		// Most packets are rather short, so allocating 16 bytes is more than enough at first.
		w.buf = make([]byte, 3, 16)
	}
}

func (w *Writer) SetId(id byte) *Writer {
	w.setupBuffer()
	w.buf[2] = byte(id)
	return w
}

func (w *Writer) WriteBytes(v []byte) *Writer {
	w.setupBuffer()
	w.buf = append(w.buf, v...)
	return w
}

func (w *Writer) WriteByte(v byte) *Writer {
	w.setupBuffer()
	w.buf = append(w.buf, v)
	return w
}

func (w *Writer) WriteBool(v bool) *Writer {
	if v {
		return w.WriteByte(1)
	} else {
		return w.WriteByte(0)
	}
}

func (w *Writer) WriteUint16(v uint16) *Writer {
	w.setupBuffer()
	vbuf := make([]byte, 2)
	binary.LittleEndian.PutUint16(vbuf, v)
	w.buf = append(w.buf, vbuf...)
	return w
}

func (w *Writer) WriteInt16(v int16) *Writer {
	return w.WriteUint16(uint16(v))
}

func (w *Writer) WriteUint32(v uint32) *Writer {
	w.setupBuffer()
	vbuf := make([]byte, 4)
	binary.LittleEndian.PutUint32(vbuf, v)
	w.buf = append(w.buf, vbuf...)
	return w
}

func (w *Writer) WriteInt32(v int32) *Writer {
	return w.WriteUint32(uint32(v))
}

func (w *Writer) WriteUint64(v uint64) *Writer {
	w.setupBuffer()
	vbuf := make([]byte, 8)
	binary.LittleEndian.PutUint64(vbuf, v)
	w.buf = append(w.buf, vbuf...)
	return w
}

func (w *Writer) WriteInt64(v int64) *Writer {
	return w.WriteUint64(uint64(v))
}

func (w *Writer) WriteFloat32(v float32) *Writer {
	return w.WriteUint32(math.Float32bits(v))
}

func (w *Writer) WriteFloat64(v float64) *Writer {
	return w.WriteUint64(math.Float64bits(v))
}

func (w *Writer) WriteString(v string) *Writer {
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

// A utility function arout w.WriteString(fmt.Sprintf(format, args...))
func (w *Writer) WriteStringf(format string, a ...interface{}) *Writer {
	return w.WriteString(fmt.Sprintf(format, a...))
}
