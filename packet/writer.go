package packet

import (
	"encoding/binary"
)

// W is used to efficiently build a terraria tcp packet
// Do not copy a non-zero Writer.
type Writer struct {
	buf []byte
}

func (pw *Writer) setupBuffer() {
	if len(pw.buf) < 3 {
		// Most packets are rather short, so allocating 16 bytes is more than enough at first.
		pw.buf = make([]byte, 3, 16)
	}
}

func (pw *Writer) Packet() *Packet {
	binary.LittleEndian.PutUint16(pw.buf, uint16(len(pw.buf)))
	buf := make([]byte, len(pw.buf))
	copy(buf, pw.buf)
	return &Packet{ptr: 3, buf: buf}
}

func (pw *Writer) SetType(t PacketType) *Writer {
	pw.setupBuffer()
	pw.buf[2] = byte(t)
	return pw
}

func (pw *Writer) WriteBytes(v []byte) *Writer {
	pw.setupBuffer()
	pw.buf = append(pw.buf, v...)
	return pw
}

func (pw *Writer) WriteByte(v byte) *Writer {
	pw.setupBuffer()
	pw.buf = append(pw.buf, v)
	return pw
}

func (pw *Writer) WriteUint16(v uint16) *Writer {
	pw.setupBuffer()
	vbuf := make([]byte, 2)
	binary.LittleEndian.PutUint16(vbuf, v)
	pw.buf = append(pw.buf, vbuf...)
	return pw
}

func (pw *Writer) WriteInt16(v int16) *Writer {
	return pw.WriteUint16(uint16(v))
}

func (pw *Writer) WriteUint32(v uint32) *Writer {
	pw.setupBuffer()
	vbuf := make([]byte, 4)
	binary.LittleEndian.PutUint32(vbuf, v)
	pw.buf = append(pw.buf, vbuf...)
	return pw
}

func (pw *Writer) WriteInt32(v int32) *Writer {
	return pw.WriteUint32(uint32(v))
}

func (pw *Writer) WriteUint64(v uint64) *Writer {
	pw.setupBuffer()
	vbuf := make([]byte, 8)
	binary.LittleEndian.PutUint64(vbuf, v)
	pw.buf = append(pw.buf, vbuf...)
	return pw
}

func (pw *Writer) PutInt64(v int64) *Writer {
	return pw.WriteUint64(uint64(v))
}

func (pw *Writer) WriteString(v string) *Writer {
	pw.setupBuffer()
	if l := len(v); l >= 128 {
		pw.WriteByte(byte((l % 128) + 128))
		pw.WriteByte(byte(l / 128))
	} else {
		pw.WriteByte(byte(l))
	}
	pw.buf = append(pw.buf, v...)
	return pw
}
