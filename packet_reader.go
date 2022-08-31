package trelay

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func NewReader(buf []byte) *Reader {
	return &Reader{ptr: 3, buf: buf}
}

type Reader struct {
	ptr uint16
	buf []byte
}

func (r *Reader) String() string {
	return fmt.Sprintf("Packet{id:%d, len:%d, data:%#v}", r.buf[2], binary.LittleEndian.Uint16(r.buf[0:2]), r.buf)
}

func (r *Reader) canReadN(l uint16) bool {
	return r.ptr+l <= r.Length()
}

func (r *Reader) Length() uint16 {
	return binary.LittleEndian.Uint16(r.buf[0:2])
}

func (r *Reader) Id() byte {
	return r.buf[2]
}

// Clones packets internal buffer and returns it
func (r *Reader) Data() []byte {
	l := binary.LittleEndian.Uint16(r.buf[0:2])
	buf := make([]byte, l)
	copy(buf, r.buf[0:l])
	return buf
}

// Resets head to 3 (start of packet body)
func (r *Reader) ResetHead() {
	r.ptr = 3
}

// Advances head l bytes, returns io.EOF if unsuccessful
func (r *Reader) AdvanceHead(l uint16) error {
	if !r.canReadN(l) {
		return io.EOF
	}
	r.ptr += l
	return nil
}

// Reads and returns a byte, error is io.EOF if unsuccessful
func (r *Reader) ReadByte() (byte, error) {
	if !r.canReadN(1) {
		return 0, io.EOF
	}
	v := r.buf[r.ptr]
	r.ptr += 1
	return v, nil
}

// Reads and returns a byte, panics if unsuccessful
func (r *Reader) MustReadByte() byte { return must(r.ReadByte()) }

// Reads a byte and returns true if is is != 0, error is io.EOF if unsuccessful
func (r *Reader) ReadBool() (bool, error) {
	v, err := r.ReadByte()
	if err != nil {
		return false, err
	}
	return v != 0, nil
}

// Reads a byte and returns true if is is != 0, panics if unsuccessful
func (r *Reader) MustReadBool() bool { return must(r.ReadBool()) }

// Reads and returns l bytes, error is io.EOF if unsuccessful. Head is not advanced on error, so it is still possible to read a smaller value from packet
func (r *Reader) ReadBytes(l uint16) ([]byte, error) {
	if !r.canReadN(l) {
		return nil, io.EOF
	}
	buf := make([]byte, l)
	copy(buf, r.buf[r.ptr:r.ptr+l])
	r.ptr += l
	return buf, nil
}

// Reads and returns l bytes, panics if unsuccessful
func (r *Reader) MustReadBytes(l uint16) []byte { return must(r.ReadBytes(l)) }

// Reads and returns uint16, error is io.EOF if unsuccessful. Head is not advanced on error, so it is still possible to read a smaller value from packet
func (r *Reader) ReadUint16() (uint16, error) {
	if !r.canReadN(2) {
		return 0, io.EOF
	}
	v := binary.LittleEndian.Uint16(r.buf[r.ptr : r.ptr+2])
	r.ptr += 2
	return v, nil
}

// Reads and returns uint16, panics if unsuccessful
func (r *Reader) MustReadUint16() uint16 { return must(r.ReadUint16()) }

// Reads and returns int16, error is io.EOF if unsuccessful. Head is not advanced on error, so it is still possible to read a smaller value from packet
func (r *Reader) ReadInt16() (int16, error) {
	v, err := r.ReadUint16()
	return int16(v), err
}

// Reads and returns int16, panics if unsuccessful
func (r *Reader) MustReadInt16() int16 { return must(r.ReadInt16()) }

// Reads and returns uint32, error is io.EOF if unsuccessful. Head is not advanced on error, so it is still possible to read a smaller value from packet
func (r *Reader) ReadUint32() (uint32, error) {
	if !r.canReadN(4) {
		return 0, io.EOF
	}
	v := binary.LittleEndian.Uint32(r.buf[r.ptr : r.ptr+4])
	r.ptr += 4
	return v, nil
}

// Reads and returns uint32, panics if unsuccessful
func (r *Reader) MustReadUint32() uint32 { return must(r.ReadUint32()) }

// Reads and returns int32, error is io.EOF if unsuccessful. Head is not advanced on error, so it is still possible to read a smaller value from packet
func (r *Reader) ReadInt32() (int32, error) {
	v, err := r.ReadUint32()
	return int32(v), err
}

// Reads and returns int32, panics if unsuccessful
func (r *Reader) MustReadInt32() int32 { return must(r.ReadInt32()) }

// Reads and returns uint64, error is io.EOF if unsuccessful. Head is not advanced on error, so it is still possible to read a smaller value from packet
func (r *Reader) ReadUint64() (uint64, error) {
	if !r.canReadN(8) {
		return 0, io.EOF
	}
	v := binary.LittleEndian.Uint64(r.buf[r.ptr : r.ptr+8])
	r.ptr += 8
	return v, nil
}

// Reads and returns uint64, panics if unsuccessful
func (r *Reader) MustReadUint64() uint64 { return must(r.ReadUint64()) }

// Reads and returns int64, error is io.EOF if unsuccessful. Head is not advanced on error, so it is still possible to read a smaller value from packet
func (r *Reader) ReadInt64() (int64, error) {
	v, err := r.ReadUint64()
	return int64(v), err
}

// Reads and returns int64, panics if unsuccessful
func (r *Reader) MustReadInt64() int64 { return must(r.ReadInt64()) }

// Reads and returns float32, error is io.EOF if unsuccessful. Head is not advanced on error, so it is still possible to read a smaller value from packet
func (r *Reader) ReadFloat32() (float32, error) {
	v, err := r.ReadUint32()
	return math.Float32frombits(v), err
}

// Reads and returns a float32, panics if unsuccessful
func (r *Reader) MustReadFloat32() float32 { return must(r.ReadFloat32()) }

// Reads and returns float64, error is io.EOF if unsuccessful. Head is not advanced on error, so it is still possible to read a smaller value from packet
func (r *Reader) ReadFloat64() (float64, error) {
	v, err := r.ReadUint64()
	return math.Float64frombits(v), err
}

// Reads and returns a float64, panics if unsuccessful
func (r *Reader) MustReadFloat64() float64 { return must(r.ReadFloat64()) }

// Reads and returns a string, error is io.EOF if unsuccessful. Head is not advanced on error, so it is still possible to read a smaller value from packet
func (r *Reader) ReadString() (string, error) {
	// saves current head so it can be restored on error
	prevptr := r.ptr

	len1, err := r.ReadByte()
	if err != nil {
		return "", err
	}
	len := uint16(len1)

	if len1 >= 128 {
		len2, err := r.ReadByte()
		if err != nil {
			r.ptr = prevptr
			return "", err
		}
		len = (len - 128) + uint16(len2<<7) // I have no idea what it does, stolen from popstarfreas/Dimensions
	}

	if !r.canReadN(len) {
		r.ptr = prevptr
		return "", io.EOF
	}

	v := string(r.buf[r.ptr : r.ptr+len])
	r.ptr += len
	return v, nil
}

// Reads and returns a string, panics if unsuccessful
func (r *Reader) MustReadString() string { return must(r.ReadString()) }

func ReadPacket(r io.Reader) (Packet, error) {
	head := make([]byte, 3)
	if n, err := r.Read(head); err != nil {
		return nil, err
	} else if n < 3 {
		// todo:
		// Technically, it is possible to read only part of the head if client's internet
		// is very slow. Shouldn't be an error though, but I (@btvoidx) am too lazy to
		// properly implement this part in the loop below, so this can be safely marked as todo
		return nil, io.EOF
	}

	ln := binary.LittleEndian.Uint16(head[0:2])

	p := make(basicPacket, ln)
	copy(p, head)

	// Read data, if any
	ptr := uint16(3)
	for ptr < ln {
		br, err := r.Read(p[ptr:ln])
		if err != nil {
			return nil, err
		}
		ptr += uint16(br)
	}

	return p, nil
}
