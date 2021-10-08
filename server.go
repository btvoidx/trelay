package trelay

import (
	"errors"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

type Server interface {
	Id() int
	Start() error
	Stop() error
	Addr() string
	RemoteAddr() string
	SetLogger(log log.FieldLogger) Server
	LoadPlugin(p Plugin) Server
	LoadPlugins(p []Plugin) Server
}

type server struct {
	id      int
	l       net.Listener
	log     log.FieldLogger
	running bool
	plugins []Plugin
	addr    string
	raddr   string
}

func NewServer(address string, remoteadress string) Server {
	return &server{
		id:      time.Now().Nanosecond(),
		log:     log.StandardLogger(),
		addr:    address,
		raddr:   remoteadress,
		plugins: make([]Plugin, 0),
	}
}

func (s *server) Id() int { return s.id }

func (s *server) Start() (err error) {
	s.l, err = net.Listen("tcp4", s.addr)
	if err != nil {
		return err
	}

	s.running = true

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
				"SessionId": session.Id(),
			}).Infof("Session with %s was opened", session.ClientConn().RemoteAddr())

			for _, plugin := range s.plugins {
				plugin.OnSessionOpen(session)
			}

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

	for _, plugin := range s.plugins {
		plugin.OnServerStop(s)
	}

	s.running = false
	return s.l.Close()
}

func (s *server) Addr() string       { return s.addr }
func (s *server) RemoteAddr() string { return s.raddr }

func (s *server) SetLogger(log log.FieldLogger) Server {
	s.log = log
	return s
}

// Load plugin into a server.
// Multiple servers can use on instance of a plugin.
// To forcefully prevent this plugin's OnLoad method should return a unique copy of the plugin.
// This method is not goroutine-safe.
func (s *server) LoadPlugin(p Plugin) Server {
	if s.running {
		s.log.Errorf("Failed to load plugin %s: server is running", p.Name())
		return s
	}

	p.OnServerStart(s)
	s.plugins = append(s.plugins, p)
	s.log.Infof("Loaded plugin %s", p.Name())

	return s
}

// Load multiple plugins into a server.
// Multiple servers can use on instance of a plugin.
// To forcefully prevent this plugin's OnLoad method should return a unique copy of the plugin.
// This method is not goroutine-safe.
func (s *server) LoadPlugins(p []Plugin) Server {
	for _, plugin := range p {
		s.LoadPlugin(plugin)
	}
	return s
}

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

			handled := false
			for _, plugin := range s.plugins {
				p.ResetHead()
				handled := plugin.OnClientPacket(p.Id(), p, session)
				handled := plugin.OnClientPacket(p.Type(), p, session)

				if handled {
					break
				}
			}

			sc := session.ServerConn()
			if !handled && sc != nil && !sc.Closed() {
				sc.Write(p) //nolint:errcheck
			}
		}

		for _, plugin := range s.plugins {
			plugin.OnSessionClose(session)
		}

		s.log.WithFields(log.Fields{
			"SessionId": session.Id(),
		}).Infof("Session with %s was closed", session.ClientConn().RemoteAddr())
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

			handled := false
			for _, plugin := range s.plugins {
				p.ResetHead()
				handled := plugin.OnServerPacket(p.Type(), p, session)

				if handled {
					break
				}
			}

			cc := session.ClientConn()
			if !handled && cc != nil && !cc.Closed() {
				cc.Write(p) //nolint:errcheck
			}
		}
	}()
}
