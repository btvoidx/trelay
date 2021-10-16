package trelay

import "encoding/binary"

type PacketWriter struct {
	ptr uint16
	buf [65535]byte
}

func (p *PacketWriter) ensurePointerIsAccurate() {
	if p.ptr < 3 {
		p.ptr = 3
	}
}

func (pw *PacketWriter) Packet() *Packet {
	binary.LittleEndian.PutUint16(pw.buf[0:2], pw.ptr)

	return &Packet{ptr: 3, buf: pw.buf}
}

func (pw *PacketWriter) SetType(t PacketType) *PacketWriter {
	pw.ensurePointerIsAccurate()
	pw.buf[2] = byte(t)
	return pw
}

func (pw *PacketWriter) PutByte(v byte) *PacketWriter {
	pw.ensurePointerIsAccurate()
	pw.buf[pw.ptr] = v
	pw.ptr += 1

	return pw
}

func (pw *PacketWriter) PutUint16(v uint16) *PacketWriter {
	pw.ensurePointerIsAccurate()
	binary.LittleEndian.PutUint16(pw.buf[pw.ptr:pw.ptr+2], v)
	pw.ptr += 2
	return pw
}

func (pw *PacketWriter) PutInt16(v int16) *PacketWriter {
	return pw.PutUint16(uint16(v))
}

func (pw *PacketWriter) PutUint32(v uint32) *PacketWriter {
	pw.ensurePointerIsAccurate()
	binary.LittleEndian.PutUint32(pw.buf[pw.ptr:pw.ptr+4], v)
	pw.ptr += 4
	return pw
}

func (pw *PacketWriter) PutInt32(v int32) *PacketWriter {
	return pw.PutUint32(uint32(v))
}

func (pw *PacketWriter) PutUint64(v uint64) *PacketWriter {
	pw.ensurePointerIsAccurate()
	binary.LittleEndian.PutUint64(pw.buf[pw.ptr:pw.ptr+8], v)
	pw.ptr += 8
	return pw
}

func (pw *PacketWriter) PutInt64(v int64) *PacketWriter {
	return pw.PutUint64(uint64(v))
}

func (pw *PacketWriter) PutString(v string) *PacketWriter {
	pw.ensurePointerIsAccurate()
	strlen := len(v)

	if strlen >= 128 {
		pw.PutByte(byte((strlen % 128) + 128))
		pw.PutByte(byte(strlen / 128))
	} else {
		pw.PutByte(byte(strlen))
	}

	copy(pw.buf[pw.ptr:], []byte(v))
	pw.ptr += uint16(strlen)

	return pw
}
