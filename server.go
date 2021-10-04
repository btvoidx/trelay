package trelay

import (
	"errors"
	"net"

	log "github.com/sirupsen/logrus"
)

type Server interface {
	Start() error
	Stop() error
	Addr() string
	RemoteAddr() string
	SetLogger(log log.FieldLogger) Server
}

type server struct {
	l     net.Listener
	log   log.FieldLogger
	addr  string
	raddr string
}

func NewServer(address string, remoteadress string) Server {
	return &server{
		log:   log.StandardLogger(),
		addr:  address,
		raddr: remoteadress,
	}
}

func (s *server) SetLogger(log log.FieldLogger) Server {
	s.log = log
	return s
}

func (s *server) Start() (err error) {
	s.l, err = net.Listen("tcp4", s.addr)
	if err != nil {
		return err
	}

	s.log.Infof("Server started on %s", s.addr)
	s.log.Infof("Proxying to %s", s.raddr)

	go func() {
		for {
			nc, err := s.l.Accept()
			if err != nil && errors.Is(err, net.ErrClosed) {
				break
			}

			session := NewSession()
			session.SetClientConn(NewConn(nc))

			s.log.WithFields(log.Fields{
				"session": session.Id(),
				"remote":  nc.RemoteAddr().String(),
			}).Info("Session opened")

			//todo OnSessionOpen

			if session.Closed() {
				continue
			}

			s.handleSession(session)
		}
	}()

	return nil
}

func (s *server) Stop() (err error) {
	s.log.Infof("Server stopped")
	return s.l.Close()
}

func (s *server) Addr() string       { return s.addr }
func (s *server) RemoteAddr() string { return s.raddr }

func (s *server) handleSession(session Session) {
	sc, err := net.Dial("tcp4", s.raddr)
	if err != nil {
		s.log.Errorf("Failed to connect to %s : %s", s.raddr, err.Error())
		return
	}

	session.SetServerConn(NewConn(sc))

	go func() {
		for {
			if session.Closed() {
				break
			}

			cc := session.ClientConn()
			if cc == nil {
				continue
			}
			if cc.Closed() {
				session.Close()
			}

			p, err := cc.Read()
			if err != nil {
				continue
			}

			//todo middleware here
			handled := false

			sc := session.ServerConn()
			if !handled && sc != nil && !sc.Closed() {
				sc.Write(p) //nolint:errcheck
			}
		}

		//todo OnSessionClose

		s.log.WithFields(log.Fields{
			"session": session.Id(),
			"remote":  session.ClientConn().RemoteAddr(),
		}).Info("Session closed")
	}()

	go func() {
		for {
			if session.Closed() {
				break
			}

			sc := session.ServerConn()
			if sc == nil {
				continue
			}
			if sc.Closed() {
				session.Close()
			}

			p, err := sc.Read()
			if err != nil {
				continue
			}

			//todo middleware here
			handled := false

			cc := session.ClientConn()
			if !handled && cc != nil && !cc.Closed() {
				cc.Write(p) //nolint:errcheck
			}
		}
	}()
}
