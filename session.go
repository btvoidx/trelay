package trelay

import (
	"io"
)

type PacketWriter interface {
	WritePacket(p Packet) error
	io.Writer
}

type packetwriter struct {
	io.ReadWriter
}

func (pw *packetwriter) WritePacket(p Packet) error {
	_, err := pw.Write(p.Data())
	return err
}

type Session interface {
	Client() PacketWriter
	// Remote is a default target for unhandled packets.
	// Can be nil if client is not yet connected to any server.
	Remote() PacketWriter
	SetRemote(io.ReadWriteCloser)

	Close() error
}

var _ Session = (*session)(nil)

type session struct {
	client io.ReadWriteCloser
	remote io.ReadWriteCloser
}

func (s *session) Client() PacketWriter {
	return &packetwriter{s.client}
}

func (s *session) Remote() PacketWriter {
	return &packetwriter{s.remote}
}

func (s *session) SetRemote(r io.ReadWriteCloser) {
	s.remote = r
}

func (s *session) Close() error {
	s.client.Close()
	if s.remote != nil {
		s.remote.Close()
	}
	return nil
}
