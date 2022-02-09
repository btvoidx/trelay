package trelay

import (
	"net"
)

// net.Conn wrapper to allow easier use with packets.
//
// Underlying Conn in PacketConn should never be nil
type PacketConn struct {
	net.Conn
}

func (pc PacketConn) Read() (*Packet, error) {
	return ReadPacket(pc.Conn)
}

func (pc PacketConn) Write(p *Packet) error {
	_, err := pc.Conn.Write(p.Data())
	return err
}
