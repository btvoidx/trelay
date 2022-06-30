package packet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriterSetType(t *testing.T) {
	assert := assert.New(t)
	p := (&Writer{}).SetType(ConnectRequest).Packet()

	assert.Equal([]byte{3, 0, 1}, p.Data())
}

func TestWriterWriteByte(t *testing.T) {
	assert := assert.New(t)
	p := (&Writer{}).WriteByte(1).Packet()

	assert.Equal([]byte{4, 0, 0, 1}, p.Data())
}

func TestWriterWriteBytes(t *testing.T) {
	assert := assert.New(t)
	p := (&Writer{}).WriteBytes([]byte{0, 255, 19, 81}).Packet()

	assert.Equal([]byte{7, 0, 0, 0, 255, 19, 81}, p.Data())
}

func TestWriterWriteUint16(t *testing.T) {
	assert := assert.New(t)
	p := (&Writer{}).WriteUint16(22444).Packet()

	assert.Equal([]byte{5, 0, 0, 172, 87}, p.Data())
}

func TestWriterWriteInt16(t *testing.T) {
	assert := assert.New(t)
	p := (&Writer{}).WriteInt16(-139).Packet()

	assert.Equal([]byte{5, 0, 0, 117, 255}, p.Data())
}
