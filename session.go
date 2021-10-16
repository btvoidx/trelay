package trelay

// Holds two connections and allows to close both of them.
type Session interface {
	// Random int in range [0,1000000000)
	Id() int
	ClientConn() PacketConn
	SetClientConn(conn PacketConn)

	ServerConn() PacketConn
	SetServerConn(conn PacketConn)

	// Closes both ClientConn and ServerConn.
	Close()

	// Returns true if both ClientConn and ServerConn are closed, otherwise false.
	Closed() bool
}

type session struct {
	id int
	cc PacketConn
	sc PacketConn
}

// Creates new empty Session.
func NewSession() Session {
	return &session{
		id: rand.Intn(1000000000),
	}
}

func (s *session) Id() int { return s.id }

func (s *session) ClientConn() PacketConn        { return s.cc }
func (s *session) SetClientConn(conn PacketConn) { s.cc = conn }

func (s *session) ServerConn() PacketConn        { return s.sc }
func (s *session) SetServerConn(conn PacketConn) { s.sc = conn }

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
