package trelay

import (
	"errors"
	"net"
	"sync"
	"sync/atomic"
)

type Handler interface {
	// Called when a new client sends their first packet.
	ClientConnect(Session)
	// Called when a client disconnects.
	ClientDisconnect(Session)
	// Called when client sends a packet.
	// Returning `true` will prevent the packet from being sent to server.
	ClientPacket(Session, Packet) (block bool)
	// Called when remote server sends a packet.
	// Returning `true` will prevent the packet from being sent to client.
	RemotePacket(Session, Packet) (block bool)
}

// Creates a default server and starts it.
func ListenAndServe(addr string, h Handler) error {
	return (&Server{Addr: addr, Handler: h}).ListenAndServe()
}

type Server struct {
	l net.Listener

	// The address to listen on.
	Addr string

	// Replacing the handler will not affect
	// established connections.
	Handler Handler
}

// Starts the server
func (s *Server) ListenAndServe() (err error) {
	s.l, err = net.Listen("tcp4", s.Addr)
	if err != nil {
		return err
	}

	for {
		nc, err := s.l.Accept()
		if err != nil && errors.Is(err, net.ErrClosed) {
			return nil
		}

		go s.handleSession(&session{
			client: nc,
		})
	}
}

func (s *Server) Stop() (err error) {
	return s.l.Close()
}

func (s *Server) handleSession(session *session) {
	onConnectCalled := false
	onConnect := s.Handler.ClientConnect
	onDisconnect := s.Handler.ClientDisconnect
	onClientPacket := s.Handler.ClientPacket
	onServerPacket := s.Handler.RemotePacket

	stopped := atomic.Bool{}
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		for !stopped.Load() {
			p, err := ReadPacket(session.client)
			if err != nil {
				break
			}

			if !onConnectCalled {
				onConnect(session)
				onConnectCalled = true
			}

			if !onClientPacket(session, p) && session.remote != nil {
				session.Remote().WritePacket(p)
			}
		}

		session.client.Close()
		stopped.Store(true)
		wg.Done()
	}()

	go func() {
		for !stopped.Load() {
			if session.remote == nil {
				continue
			}

			p, err := ReadPacket(session.remote)
			// if err == io.EOF {
			if err != nil {
				break
			}

			if !onServerPacket(session, p) {
				session.Client().WritePacket(p)
			}
		}

		if session.remote != nil {
			session.remote.Close()
		}

		stopped.Store(true)
		wg.Done()
	}()

	wg.Wait()

	onDisconnect(session)
}
