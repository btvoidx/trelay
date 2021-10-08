package trelay

import (
	"io"
	"net"
)

type Conn interface {
	Read() (Packet, error)
	Write(Packet) (int, error)
	RemoteAddr() string
	Close() error
	Closed() bool
}

type conn struct {
	nc     net.Conn
	closed bool
}

func NewConn(nc net.Conn) Conn {
	conn := &conn{
		nc:     nc,
		closed: false,
	}

	return conn
}

func (c *conn) Read() (Packet, error) {
	p, err := NewPacketFromReader(c.nc)

	if err != nil && err == io.EOF {
		c.Close()
		return nil, err
	}

	return p, err
}

func (c *conn) Write(p Packet) (int, error) { return c.nc.Write(p.Data()) }

func (c *conn) RemoteAddr() string { return c.nc.RemoteAddr().String() }

func (c *conn) Close() error {
	err := c.nc.Close()
	if err == nil {
		c.closed = true
	}
	return err
}

func (c *conn) Closed() bool { return c.closed }
