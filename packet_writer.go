package trelay

import "encoding/binary"

type PacketWriter struct {
	ptr uint16
	buf [65535]byte
}

func (p *PacketWriter) ensureAccuratePointer() {
	if p.ptr < 3 {
		p.ptr = 3
	}
}

func (pw *PacketWriter) Packet() *Packet {
	binary.LittleEndian.PutUint16(pw.buf[0:2], pw.ptr)
	return &Packet{ptr: 3, buf: pw.buf[0:pw.ptr]}
}

func (pw *PacketWriter) SetType(t PacketType) *PacketWriter {
	pw.ensureAccuratePointer()
	pw.buf[2] = byte(t)
	return pw
}

func (pw *PacketWriter) WriteBytes(v []byte) *PacketWriter {
	pw.ensureAccuratePointer()
	copy(pw.buf[pw.ptr:], v)
	pw.ptr += uint16(len(v))

	return pw
}

func (pw *PacketWriter) WriteByte(v byte) *PacketWriter {
	pw.ensureAccuratePointer()
	pw.buf[pw.ptr] = v
	pw.ptr += 1

	return pw
}

func (pw *PacketWriter) WriteUint16(v uint16) *PacketWriter {
	pw.ensureAccuratePointer()
	binary.LittleEndian.PutUint16(pw.buf[pw.ptr:pw.ptr+2], v)
	pw.ptr += 2
	return pw
}

func (pw *PacketWriter) WriteInt16(v int16) *PacketWriter {
	return pw.WriteUint16(uint16(v))
}

func (pw *PacketWriter) WriteUint32(v uint32) *PacketWriter {
	pw.ensureAccuratePointer()
	binary.LittleEndian.PutUint32(pw.buf[pw.ptr:pw.ptr+4], v)
	pw.ptr += 4
	return pw
}

func (pw *PacketWriter) WriteInt32(v int32) *PacketWriter {
	return pw.WriteUint32(uint32(v))
}

func (pw *PacketWriter) WriteUint64(v uint64) *PacketWriter {
	pw.ensureAccuratePointer()
	binary.LittleEndian.PutUint64(pw.buf[pw.ptr:pw.ptr+8], v)
	pw.ptr += 8
	return pw
}

func (pw *PacketWriter) PutInt64(v int64) *PacketWriter {
	return pw.WriteUint64(uint64(v))
}

func (pw *PacketWriter) WriteString(v string) *PacketWriter {
	pw.ensureAccuratePointer()
	strlen := len(v)

	if strlen >= 128 {
		pw.WriteByte(byte((strlen % 128) + 128))
		pw.WriteByte(byte(strlen / 128))
	} else {
		pw.WriteByte(byte(strlen))
	}

	copy(pw.buf[pw.ptr:], []byte(v))
	pw.ptr += uint16(strlen)

	return pw
}
