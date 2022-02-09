package trelay

type PacketContext struct {
	packet  *Packet
	handled bool
	session *Session
}

// Returns the received packet and resets its head
func (pc *PacketContext) Packet() *Packet {
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

// Client connection
func (pc *PacketContext) Session() *Session {
	return pc.session
}
