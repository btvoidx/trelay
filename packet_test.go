package trelay

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBadLengthPacket(t *testing.T) {
	assert := assert.New(t)
	r := bytes.NewReader([]byte{2, 0})
	_, err := ReadPacket(r)

	assert.EqualError(err, "bad packet length")
}

func TestReadPacketEOF(t *testing.T) {
	assert := assert.New(t)
	r := bytes.NewReader([]byte{10, 0, 1, 0, 0, 0})
	_, err := ReadPacket(r)

	assert.EqualError(err, io.EOF.Error())
}

func TestCorrectPacketLength(t *testing.T) {
	assert := assert.New(t)
	r := bytes.NewReader([]byte{6, 0, 1, 0, 0, 1})
	p, err := ReadPacket(r)

	if !assert.NoError(err) {
		return
	}

	assert.Equal(uint16(6), p.Length())
}

func TestPacketResetHead(t *testing.T) {
	assert := assert.New(t)
	r := bytes.NewReader([]byte{6, 0, 1, 0, 5, 1})
	p, err := ReadPacket(r)

	if !assert.NoError(err) {
		return
	}

	assert.Equal(uint16(3), p.ptr)
	p.MustReadUint16()
	assert.Equal(uint16(5), p.ptr)
	p.ResetHead()
	assert.Equal(uint16(3), p.ptr)
}

func TestCorrectPacketType(t *testing.T) {
	assert := assert.New(t)
	r := bytes.NewReader([]byte{3, 0, byte(ConnectRequest)}) // Valid packet
	p, err := ReadPacket(r)

	if assert.NoError(err) {
		assert.Equal(ConnectRequest, p.Type())
	}
}

func TestReadByte(t *testing.T) {
	assert := assert.New(t)
	r := bytes.NewReader([]byte{4, 0, 0, 5})
	p, err := ReadPacket(r)

	if !assert.NoError(err) {
		return
	}

	v, err := p.ReadByte()
	if !assert.NoError(err) {
		return
	}

	assert.Equal(byte(5), v)
}

func TestReadBytes(t *testing.T) {
	assert := assert.New(t)
	r := bytes.NewReader([]byte{6, 0, 0, 5, 6, 7})
	p, err := ReadPacket(r)

	if !assert.NoError(err) {
		return
	}

	buf, err := p.ReadBytes(3)
	if !assert.NoError(err) {
		return
	}

	assert.Equal(byte(5), buf[0])
	assert.Equal(byte(6), buf[1])
	assert.Equal(byte(7), buf[2])
}

func TestReadUint16(t *testing.T) {
	assert := assert.New(t)
	r := bytes.NewReader([]byte{5, 0, 1, 172, 87})
	p, err := ReadPacket(r)

	if !assert.NoError(err) {
		return
	}

	v, err := p.ReadUint16()
	if !assert.NoError(err) {
		return
	}

	assert.Equal(uint16(22444), v)
}

func TestReadInt16(t *testing.T) {
	assert := assert.New(t)
	r := bytes.NewReader([]byte{5, 0, 1, 117, 255})
	p, err := ReadPacket(r)

	if !assert.NoError(err) {
		return
	}

	v, err := p.ReadInt16()
	if assert.NoError(err) {
		return
	}

	assert.Equal(int16(-139), v)
}

func TestReadString(t *testing.T) {
	assert := assert.New(t)
	// Encoded string: fdec2c95-f203-47b4-8256-3c3b156b251e
	r := bytes.NewReader([]byte{40, 0, 68, 36, 102, 100, 101, 99, 50, 99, 57, 53, 45, 102, 50, 48, 51, 45, 52, 55, 98, 52, 45, 56, 50, 53, 54, 45, 51, 99, 51, 98, 49, 53, 54, 98, 50, 53, 49, 101})
	p, err := ReadPacket(r)

	if !assert.NoError(err) {
		return
	}

	v, err := p.ReadString()
	if !assert.NoError(err) {
		return
	}

	assert.Equal("fdec2c95-f203-47b4-8256-3c3b156b251e", v)
}
