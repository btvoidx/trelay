package trelay

import "time"

// Holds two connections and allows to close both of them
type Session interface {
	Id() int
	ClientConn() Conn
	SetClientConn(conn Conn)
	ServerConn() Conn
	SetServerConn(conn Conn)
	Close()
	Closed() bool
}

type session struct {
	id int
	cc Conn
	sc Conn
}

func NewSession() Session {
	return &session{
		id: time.Now().Nanosecond(),
	}
}

func (s *session) Id() int { return s.id }

func (s *session) ClientConn() Conn        { return s.cc }
func (s *session) SetClientConn(conn Conn) { s.cc = conn }

func (s *session) ServerConn() Conn        { return s.sc }
func (s *session) SetServerConn(conn Conn) { s.sc = conn }

func (s *session) Close() {
	if s.cc != nil && !s.cc.Closed() {
		s.cc.Close()
	}

	if s.sc != nil && !s.sc.Closed() {
		s.sc.Close()
	}
}

func (s *session) Closed() bool {
	return (s.cc == nil || s.cc.Closed()) && (s.sc == nil || s.sc.Closed())
}
