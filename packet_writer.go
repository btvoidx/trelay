package trelay

import "encoding/binary"

// Packet writer easy creation of custom tcp packets
type PacketWriter struct {
	buf []byte
}

func (pw *PacketWriter) setupBuffer() {
	if len(pw.buf) < 3 {
		// 128 bytes is about right of a tradeoff between allocating more and allocating multiple times for shorter and longer packets
		pw.buf = make([]byte, 3, 128)
	}
}

func (pw *PacketWriter) Packet() *Packet {
	binary.LittleEndian.PutUint16(pw.buf, uint16(len(pw.buf)))
	buf := make([]byte, len(pw.buf))
	copy(buf, pw.buf)
	return &Packet{ptr: 3, buf: buf}
}

func (pw *PacketWriter) SetType(t PacketType) *PacketWriter {
	pw.setupBuffer()
	pw.buf[2] = byte(t)
	return pw
}

func (pw *PacketWriter) WriteBytes(v []byte) *PacketWriter {
	pw.setupBuffer()
	pw.buf = append(pw.buf, v...)
	return pw
}

func (pw *PacketWriter) WriteByte(v byte) *PacketWriter {
	pw.setupBuffer()
	pw.buf = append(pw.buf, v)
	return pw
}

func (pw *PacketWriter) WriteUint16(v uint16) *PacketWriter {
	pw.setupBuffer()
	vbuf := make([]byte, 2)
	binary.LittleEndian.PutUint16(vbuf, v)
	pw.buf = append(pw.buf, vbuf...)
	return pw
}

func (pw *PacketWriter) WriteInt16(v int16) *PacketWriter {
	return pw.WriteUint16(uint16(v))
}

func (pw *PacketWriter) WriteUint32(v uint32) *PacketWriter {
	pw.setupBuffer()
	vbuf := make([]byte, 4)
	binary.LittleEndian.PutUint32(vbuf, v)
	pw.buf = append(pw.buf, vbuf...)
	return pw
}

func (pw *PacketWriter) WriteInt32(v int32) *PacketWriter {
	return pw.WriteUint32(uint32(v))
}

func (pw *PacketWriter) WriteUint64(v uint64) *PacketWriter {
	pw.setupBuffer()
	vbuf := make([]byte, 8)
	binary.LittleEndian.PutUint64(vbuf, v)
	pw.buf = append(pw.buf, vbuf...)
	return pw
}

func (pw *PacketWriter) PutInt64(v int64) *PacketWriter {
	return pw.WriteUint64(uint64(v))
}

func (pw *PacketWriter) WriteString(v string) *PacketWriter {
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
