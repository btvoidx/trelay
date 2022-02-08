package trelay

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPacketSetType(t *testing.T) {
	assert := assert.New(t)
	p := (&PacketWriter{}).SetType(ConnectRequest).Packet()

	assert.Equal([]byte{3, 0, 1}, p.Data())
}

func TestPacketWriteByte(t *testing.T) {
	assert := assert.New(t)
	p := (&PacketWriter{}).WriteByte(1).Packet()

	assert.Equal([]byte{4, 0, 0, 1}, p.Data())
}

func TestPacketWriteBytes(t *testing.T) {
	assert := assert.New(t)
	p := (&PacketWriter{}).WriteBytes([]byte{0, 255, 19, 81}).Packet()

	assert.Equal([]byte{7, 0, 0, 0, 255, 19, 81}, p.Data())
}

func TestPacketWriteUint16(t *testing.T) {
	assert := assert.New(t)
	p := (&PacketWriter{}).WriteUint16(22444).Packet()

	assert.Equal([]byte{5, 0, 0, 172, 87}, p.Data())
}

func TestPacketWriteInt16(t *testing.T) {
	assert := assert.New(t)
	p := (&PacketWriter{}).WriteInt16(-139).Packet()

	assert.Equal([]byte{5, 0, 0, 117, 255}, p.Data())
}
