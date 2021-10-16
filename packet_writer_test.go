package trelay_test

import (
	"testing"

	"github.com/btvoidx/trelay"
	"github.com/stretchr/testify/assert"
)

func TestPacketSetType(t *testing.T) {
	assert := assert.New(t)
	var pw trelay.PacketWriter
	p := pw.SetType(trelay.ConnectRequest).Packet()

	assert.Equal([]byte{3, 0, 1}, p.Data())
}

func TestPacketPutByte(t *testing.T) {
	assert := assert.New(t)
	var pw trelay.PacketWriter
	p := pw.PutByte(1).Packet()

	assert.Equal([]byte{4, 0, 0, 1}, p.Data())
}

func TestPacketPutUint16(t *testing.T) {
	assert := assert.New(t)
	var pw trelay.PacketWriter
	p := pw.PutUint16(22444).Packet()

	assert.Equal([]byte{5, 0, 0, 172, 87}, p.Data())
}

func TestPacketPutInt16(t *testing.T) {
	assert := assert.New(t)
	var pw trelay.PacketWriter
	p := pw.PutInt16(-139).Packet()

	assert.Equal([]byte{5, 0, 0, 117, 255}, p.Data())
}
