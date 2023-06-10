package trelay

import (
	"encoding"
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

type Packet struct {
	id  byte
	len uint16
	ptr uint16
	buf []byte
}

var (
	_ io.Reader     = (*Packet)(nil)
	_ io.ReaderFrom = (*Packet)(nil)
	// _ io.ReaderAt   = (*Packet)(nil) // TODO
	_ io.WriterTo = (*Packet)(nil)
)

func (p *Packet) Id() byte { return p.id }

// Len includes length of packet header (3 bytes)
func (p *Packet) Len() uint16 { return p.len }

func (p *Packet) Read(b []byte) (n int, err error) {
	if int(p.ptr) >= len(p.buf) {
		return 0, io.EOF
	}
	n = copy(b, p.buf[p.ptr:])
	p.ptr += uint16(n)
	return
}

// Will read an entire packet from r and reset internal buffer pointer,
// assuming first three bytes read from r form a valid packet header.
//
// As opposed to definition of io.ReaderFrom, reads only up to
// length given in the header. Will reuse internal buffer when possible.
func (p *Packet) ReadFrom(r io.Reader) (n int64, err error) {
	p.ptr = 0

	nn, err := Fscan(r, &p.len, &p.id)
	n += int64(nn)
	if err != nil {
		return n, err
	}

	if p.len <= 3 {
		// effectively empties the buffer so Reads fail
		// but preserves it to allow for later reuse
		p.buf = p.buf[:0]
		return n, nil
	}

	if int(p.len-3) > cap(p.buf) {
		p.buf = make([]byte, p.len-3)
	}

	p.buf = p.buf[0 : p.len-3] // reslice

	nn, err = Fscan(r, &p.buf)
	n += int64(nn)
	return n, err
}

func (p Packet) WriteTo(w io.Writer) (n int64, err error) {
	nn, err := Fprint(w, p.len, p.id, p.buf)
	return int64(nn), err
}

// func (p *Packet) ReadAt(b []byte, off int64) (n int, err error)

// Builder helps avoid manual bookkeeping of packet
// length and id as opposed to bytes.Buffer.
// Do not copy a non-zero Builder, as buffer will be shared.
// Zero value is ready to use.
type Builder struct {
	ID  byte
	buf []byte
}

var (
	_ io.Writer   = (*Builder)(nil)
	_ io.WriterTo = (*Builder)(nil)
)

func (b *Builder) setupBuffer() {
	if b.buf == nil {
		b.buf = make([]byte, 3, 32)
	}
}

func (b *Builder) Write(d []byte) (n int, err error) {
	b.setupBuffer()
	b.buf = append(b.buf, d...)
	return len(d), nil
}

// Reslices internal buffer to clear it for reuse.
func (b *Builder) Reset() {
	b.setupBuffer()
	b.buf = b.buf[:3]
}

// Updates length and ID in the internal buffer and returns it.
func (b *Builder) Bytes() []byte {
	b.setupBuffer()
	b.buf[2] = b.ID
	binary.LittleEndian.PutUint16(b.buf[0:2], uint16(len(b.buf)))
	return b.buf
}

func (b Builder) WriteTo(w io.Writer) (n int64, err error) {
	d := b.Bytes()
	nn, err := w.Write(d)
	return int64(nn), err
}

// fmt.Fscan-ish helper to read data from Terraria packets.
// Supported types: io.ReaderFrom, *[]byte, *byte, *[u]int[8/16/32/64], *float[32/64], string.
// Note that values implementing io.ReaderFrom should stop before EOF or error, otherwise they
// will consume the entire reader.
func Fscan(r io.Reader, ptrs ...any) (n int, err error) {
	var a int
	for _, p := range ptrs {
		switch p := p.(type) {
		default:
			return n, fmt.Errorf("unsupported type %T", p)
		case io.ReaderFrom:
			var a64 int64 // avoid shadowing err
			a64, err = p.ReadFrom(r)
			a = int(a64)
		case *[]byte:
			a, err = io.ReadFull(r, *p)
		case *byte:
			var buf [1]byte
			a, err = io.ReadFull(r, buf[:])
			*p = buf[0]
		case *int8:
			var buf [1]byte
			a, err = io.ReadFull(r, buf[:])
			*p = int8(buf[0])
		case *bool:
			var buf [1]byte
			a, err = io.ReadFull(r, buf[:])
			*p = buf[0] != 0
		case *uint16:
			var buf [2]byte
			a, err = io.ReadFull(r, buf[:])
			*p = binary.LittleEndian.Uint16(buf[:])
		case *int16:
			var buf [2]byte
			a, err = io.ReadFull(r, buf[:])
			*p = int16(binary.LittleEndian.Uint16(buf[:]))
		case *uint32:
			var buf [4]byte
			a, err = io.ReadFull(r, buf[:])
			*p = binary.LittleEndian.Uint32(buf[:])
		case *int32:
			var buf [4]byte
			a, err = io.ReadFull(r, buf[:])
			*p = int32(binary.LittleEndian.Uint32(buf[:]))
		case *uint64:
			var buf [8]byte
			a, err = io.ReadFull(r, buf[:])
			*p = binary.LittleEndian.Uint64(buf[:])
		case *int64:
			var buf [8]byte
			a, err = io.ReadFull(r, buf[:])
			*p = int64(binary.LittleEndian.Uint64(buf[:]))
		case *float32:
			var buf [4]byte
			a, err = io.ReadFull(r, buf[:])
			bits := binary.LittleEndian.Uint32(buf[:])
			*p = math.Float32frombits(bits)
		case *float64:
			var buf [8]byte
			a, err = io.ReadFull(r, buf[:])
			bits := binary.LittleEndian.Uint64(buf[:])
			*p = math.Float64frombits(bits)
		case *string:
			var lenb [1]byte
			a, err = io.ReadFull(r, lenb[:])
			if err != nil {
				break
			}
			n += a

			len := uint16(lenb[0])
			if len >= 128 {
				a, err = io.ReadFull(r, lenb[:])
				if err != nil {
					break
				}
				n += a

				// I have no idea what it does, stolen from popstarfreas/Dimensions
				len = (len - 128) + uint16(lenb[0]<<7)
			}

			buf := make([]byte, len)
			a, err = io.ReadFull(r, buf[:])
			*p = string(buf[:])
		}

		n += a
		if err != nil {
			return n, err
		}
	}

	return
}

// fmt.Fprint-ish helper to write data as Terraria expects it.
// Supported types: io.WriterTo, encoding.BinaryMarshaler,
// []byte, [u]int[8/16/32/64], float[32/64], string.
//
// Does not write packet header, it is assumed to be written earlier.
// Use *Builder to build a packet prior to sending it.
func Fprint(w io.Writer, v ...any) (n int, err error) {
	var a int
	for _, v := range v {
		switch v := v.(type) {
		default:
			return n, fmt.Errorf("unsupported type %T", v)
		case io.WriterTo:
			var a64 int64 // avoid shadowing err
			a64, err = v.WriteTo(w)
			a = int(a64)
		case encoding.BinaryMarshaler:
			var data []byte
			data, err = v.MarshalBinary()
			if err != nil {
				break
			}
			a, err = w.Write(data)
		case []byte:
			a, err = w.Write(v)
		case byte:
			buf := [1]byte{v}
			a, err = w.Write(buf[:])
		case int8:
			buf := [1]byte{byte(v)}
			a, err = w.Write(buf[:])
		case bool:
			var buf [1]byte
			if v {
				buf[0] = 1
			}
			a, err = w.Write(buf[:])
		case uint16:
			var buf [2]byte
			binary.LittleEndian.PutUint16(buf[:], v)
			a, err = w.Write(buf[:])
		case int16:
			var buf [2]byte
			binary.LittleEndian.PutUint16(buf[:], uint16(v))
			a, err = w.Write(buf[:])
		case uint32:
			var buf [4]byte
			binary.LittleEndian.PutUint32(buf[:], v)
			a, err = w.Write(buf[:])
		case int32:
			var buf [4]byte
			binary.LittleEndian.PutUint32(buf[:], uint32(v))
			a, err = w.Write(buf[:])
		case uint64:
			var buf [8]byte
			binary.LittleEndian.PutUint64(buf[:], v)
			a, err = w.Write(buf[:])
		case int64:
			var buf [8]byte
			binary.LittleEndian.PutUint64(buf[:], uint64(v))
			a, err = w.Write(buf[:])
		case string:
			if ln := len(v); ln >= 128 {
				buf := [2]byte{byte((ln % 128) + 128), byte(ln / 128)}
				a, err = w.Write(buf[:])
			} else {
				buf := [1]byte{byte(ln)}
				a, err = w.Write(buf[:])
			}

			n += a
			if err != nil {
				break
			}

			a, err = w.Write([]byte(v)) // do unsafe to avoid a heap copy?
		}

		n += a
		if err != nil {
			return n, err
		}
	}
	return
}
