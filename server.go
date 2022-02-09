package trelay

import (
	"errors"
	"net"
)

type Options struct {
	Addr                  string
	RemoteAddr            string
	DisableStandartLogger bool
}

type Server struct {
	l       net.Listener
	plugins []Plugin
	addr    string
	raddr   string
}

func NewServer(c Options) *Server {
	s := &Server{
		addr:    c.Addr,
		raddr:   c.RemoteAddr,
		plugins: make([]Plugin, 0),
	}

	return s
}

// Starts server as a goroutine.
func (s *Server) Start() (err error) {
	s.l, err = net.Listen("tcp4", s.addr)
	if err != nil {
		return err
	}

	for _, plugin := range s.plugins {
		plugin.OnServerStart()
	}

	go func() {
		for {
			nc, err := s.l.Accept()
			if err != nil && errors.Is(err, net.ErrClosed) {
				return
			}

			session := &Session{
				Client: PacketConn{nc},
			}

			s.handleSession(session)
		}
	}()

	return nil
}

// Stops server's net.Listener and calls OnServerStop on all loaded plugins.
func (s *Server) Stop() (err error) {
	for _, plugin := range s.plugins {
		plugin.OnServerStop()
	}

	return s.l.Close()
}

func (s *Server) Addr() string       { return s.addr }
func (s *Server) RemoteAddr() string { return s.raddr }

func (s *Server) LoadPlugin(loader func(*Server) Plugin) *Server {
	s.plugins = append(s.plugins, loader(s))
	return s
}

func (s *Server) handleSession(session *Session) {
	for _, plugin := range s.plugins {
		plugin.OnSessionOpen(session)
	}

	if session.Server.Conn == nil {
		sc, err := net.Dial("tcp4", s.raddr)
		if err != nil {
			return
		}
		session.Server = PacketConn{sc}
	}

	go func() {
		for {
			p, err := session.Client.Read()
			// if err == io.EOF {
			if err != nil {
				break
			}

			ctx := &PacketContext{
				packet:  p,
				session: session,
			}

			for _, plugin := range s.plugins {
				plugin.OnClientPacket(ctx)
			}

			if ctx.handled {
				continue
			}

			if session.Server.Write(p) != nil {
				break
			}
		}

		for _, plugin := range s.plugins {
			plugin.OnSessionClose(session)
		}
	}()

	go func() {
		for {
			p, err := session.Server.Read()
			if err != nil {
				break
			}

			ctx := &PacketContext{
				packet:  p,
				session: session,
			}

			for _, plugin := range s.plugins {
				plugin.OnServerPacket(ctx)
			}

			if ctx.handled {
				continue
			}

			if session.Client.Write(p) != nil {
				break
			}
		}
	}()
}
