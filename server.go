package trelay

import (
	"errors"
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
)

type Server interface {
	// Random int in range [0,1000000000).
	Id() int

	// Starts server as a goroutine.
	Start() error

	// Stops server's net.Listener and calls OnServerStop on all loaded plugins.
	Stop() error

	// Retruns server address.
	Addr() string

	// Returns remote address to which server routes by default.
	RemoteAddr() string

	Logger() log.FieldLogger
	SetLogger(log log.FieldLogger) Server

	// Load plugin into a server.
	//
	// This method is not goroutine-safe.
	LoadPlugin(loader func(Server) Plugin) Server

	// Look comment for LoadPlugin().
	LoadPlugins(loaders []func(Server) Plugin) Server
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

func NewServer(address string, remoteaddress string) Server {
	return &server{
		id:      rand.Intn(1000000000),
		log:     log.StandardLogger(),
		addr:    address,
		raddr:   remoteaddress,
		plugins: make([]Plugin, 0),
	}
}

func (s *server) Id() int { return s.id }

func (s *server) Start() (err error) {
	if s.running {
		return fmt.Errorf("failed to start server: already running")
	}

	s.l, err = net.Listen("tcp4", s.addr)
	if err != nil {
		return err
	}

	s.running = true

	s.log.Infof("Server started on \"%s\"", s.addr)
	s.log.Infof("Proxying to \"%s\"", s.raddr)

	for _, plugin := range s.plugins {
		plugin.OnServerStart()
	}

	go func() {
		for {
			nc, err := s.l.Accept()
			if err != nil && errors.Is(err, net.ErrClosed) {
				break
			}

			session := NewSession()
			session.SetClientConn(NewPacketConn(nc))

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
	if !s.running {
		return nil
	}

	s.log.Infof("Stopping server")

	for _, plugin := range s.plugins {
		plugin.OnServerStop()
	}

	err = s.l.Close()
	if err != nil {
		for _, plugin := range s.plugins {
			plugin.OnServerStart()
		}

		return err
	}

	s.running = false
	s.log.Infof("Server stopped")

	return nil
}

func (s *server) Addr() string       { return s.addr }
func (s *server) RemoteAddr() string { return s.raddr }

func (s *server) Logger() log.FieldLogger { return s.log }
func (s *server) SetLogger(log log.FieldLogger) Server {
	s.log = log
	return s
}

func (s *server) LoadPlugin(loader func(s Server) Plugin) Server {
	if s.running {
		s.log.Errorf("Failed to load plugin: server is running")
		return s
	}

	p := loader(s)
	s.plugins = append(s.plugins, p)
	s.log.Infof("Loaded plugin \"%s\"", p.Name())

	return s
}

func (s *server) LoadPlugins(loaders []func(s Server) Plugin) Server {
	for _, loader := range loaders {
		s.LoadPlugin(loader)
	}
	return s
}

func (s *server) handleSession(session Session) {
	sc, err := net.Dial("tcp4", s.raddr)
	if err != nil {
		s.log.Errorf("Failed to connect to %s : %s", s.raddr, err.Error())
		return
	}

	session.SetServerConn(NewPacketConn(sc))

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
