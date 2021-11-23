package trelay

type PacketContext interface {
	Handled() bool
	SetHandled()
	Server() Server
}

type packetcontext struct {
	handled bool
	server  Server
}

// Returns true if a plugin already marked this packet as handled
func (pc *packetcontext) Handled() bool {
	return pc.handled
}

// Marks packet as handled. Other plugins will still get it, but not the intended packet reciever
func (pc *packetcontext) SetHandled() {
	pc.handled = true
}

// Server which handles the packet
func (pc *packetcontext) Server() Server {
	return pc.server
}
