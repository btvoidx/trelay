package trelay

// Holds two connections and allows to close them or swap them out
type Session struct {
	Client PacketConn
	Server PacketConn
}

// Closes both connections
func (s Session) Close() {
	s.Client.Close()
	s.Server.Close()
}
