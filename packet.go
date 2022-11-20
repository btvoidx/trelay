package trelay

import (
	"encoding/binary"
	"fmt"
)

type Packet interface {
	// Returns the packet Id.
	Id() byte

	// Returns the packet length.
	Length() uint16

	// Returns a copy of the packet buffer as sent by Terraria (includes packet length and Id).
	Data() []byte
}

// Basic packet implementation. Eventually will be replaced by pre-parsed packets.
type basicPacket []byte

var _ Packet = (*basicPacket)(nil)

func (p basicPacket) Id() byte {
	return p[2]
}

func (p basicPacket) Length() uint16 {
	return binary.LittleEndian.Uint16(p[0:2])
}

func (p basicPacket) Data() []byte {
	buf := make([]byte, len(p))
	copy(buf, p)
	return buf
}

func (p basicPacket) String() string {
	return fmt.Sprintf("Packet{id:%d, len:%d, data:%v}", p[2], binary.LittleEndian.Uint16(p[0:2]), []byte(p))
}
