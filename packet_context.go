package trelay

import (
	"github.com/btvoidx/trelay/packet"
)

type PacketContext struct {
	packet  *packet.Packet
	handled bool
}

// Returns the received packet and resets its head
func (pc *PacketContext) Packet() *packet.Packet {
	pc.packet.ResetHead()
	return pc.packet
}

// Returns true if a plugin has marked this packet as handled
func (pc *PacketContext) Handled() bool {
	return pc.handled
}

// Marks packet as handled. Other plugins will still receive it, but it will not be sent over network
func (pc *PacketContext) SetHandled() {
	pc.handled = true
}
