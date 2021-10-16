package trelay_test

import (
	"bytes"
	"testing"

	"github.com/btvoidx/trelay"
	"github.com/stretchr/testify/assert"
)

func TestShortPacket(t *testing.T) {
	assert := assert.New(t)
	r := bytes.NewReader([]byte{2, 0})
	_, err := trelay.ReadPacket(r)

	assert.EqualError(err, "packet is too short: expecting 3 bytes, got 2")
}

func TestFalseLengthPacket(t *testing.T) {
	assert := assert.New(t)
	r := bytes.NewReader([]byte{10, 0, 1, 0, 0, 0})
	_, err := trelay.ReadPacket(r)

	assert.EqualError(err, "EOF")
}

func TestCorrectPacketLength(t *testing.T) {
	assert := assert.New(t)
	r := bytes.NewReader([]byte{6, 0, 1, 0, 0, 1})
	p, err := trelay.ReadPacket(r)

	if assert.NoError(err) {
		assert.Equal(uint16(6), p.Length())
	}
}

func TestCorrectPacketType(t *testing.T) {
	assert := assert.New(t)
	r := bytes.NewReader([]byte{3, 0, byte(trelay.ConnectRequest)}) // Valid packet
	p, err := trelay.ReadPacket(r)

	if assert.NoError(err) {
		assert.Equal(trelay.ConnectRequest, p.Type())
	}
}

func TestReadByte(t *testing.T) {
	assert := assert.New(t)
	r := bytes.NewReader([]byte{4, 0, 0, 5})
	p, err := trelay.ReadPacket(r)

	if assert.NoError(err) {
		assert.Equal(byte(5), p.ReadByte())
	}
}

func TestReadUint16(t *testing.T) {
	assert := assert.New(t)
	r := bytes.NewReader([]byte{5, 0, 1, 172, 87})
	p, err := trelay.ReadPacket(r)

	if assert.NoError(err) {
		assert.Equal(uint16(22444), p.ReadUint16())
	}
}

func TestReadInt16(t *testing.T) {
	assert := assert.New(t)
	r := bytes.NewReader([]byte{5, 0, 1, 117, 255})
	p, err := trelay.ReadPacket(r)

	if assert.NoError(err) {
		assert.Equal(int16(-139), p.ReadInt16())
	}
}

func TestReadString(t *testing.T) {
	assert := assert.New(t)
	// Encoded string: fdec2c95-f203-47b4-8256-3c3b156b251e
	r := bytes.NewReader([]byte{40, 0, 68, 36, 102, 100, 101, 99, 50, 99, 57, 53, 45, 102, 50, 48, 51, 45, 52, 55, 98, 52, 45, 56, 50, 53, 54, 45, 51, 99, 51, 98, 49, 53, 54, 98, 50, 53, 49, 101})
	p, err := trelay.ReadPacket(r)

	if assert.NoError(err) {
		assert.Equal("fdec2c95-f203-47b4-8256-3c3b156b251e", p.ReadString())
	}
}
