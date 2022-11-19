package trelay

import (
	"io"
)

type Session interface {
	Client() io.WriteCloser
	// Remote is a default target for unhandled packets.
	// Can be nil if client is not yet connected to any server.
	Remote() io.Writer
	SetRemote(io.ReadWriteCloser)
}

var _ Session = (*session)(nil)

type session struct {
	client io.ReadWriteCloser
	remote io.ReadWriteCloser
}

func (s *session) Client() io.WriteCloser {
	return s.client
}

func (s *session) Remote() io.Writer {
	return s.remote
}

func (s *session) SetRemote(r io.ReadWriteCloser) {
	s.remote = r
}
