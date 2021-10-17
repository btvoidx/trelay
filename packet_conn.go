package trelay

import (
	"io"
	"net"
)

// Wrapper around net.Conn for easy use with packets
type PacketConn interface {
	Read() (*Packet, error)
	Write(*Packet) (int, error)
	RemoteAddr() string
	Close() error
	Closed() bool
}

type pconn struct {
	nc     net.Conn
	closed bool
}

func NewPacketConn(nc net.Conn) PacketConn {
	pc := &pconn{
		nc:     nc,
		closed: false,
	}

	return pc
}

func (c *pconn) Read() (*Packet, error) {
	p, err := ReadPacket(c.nc)

	if err != nil && err == io.EOF {
		c.Close()
		return nil, err
	}

	return p, err
}

func (c *pconn) Write(p *Packet) (int, error) { return c.nc.Write(p.Data()) }

func (c *pconn) RemoteAddr() string { return c.nc.RemoteAddr().String() }

func (c *pconn) Close() error {
	err := c.nc.Close()
	if err == nil {
		c.closed = true
	}
	return err
}

func (c *pconn) Closed() bool { return c.closed }
