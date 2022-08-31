package trelay

import (
	"errors"
	"net"
	"sync"
	"sync/atomic"
)

// SessionHandler is a function that captures a session and returns two functions:
// first handles packets sent by client, second handles packets sent by Terraria server.
// Third function can be used to clean up any resources allocated that would otherwise be ignored by garbage collector.
//
// If one of the returned functions is nil, it is considered a no-op.
type SessionHandler func(Session) (fromclient func(Packet), fromserver func(Packet), cleanup func())

type Server struct {
	// The address to listen on.
	Addr string

	// The function that gets called when a new session is created.
	Handler SessionHandler

	l net.Listener
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

		session := &session{
			client: nc,
		}

		go s.handleSession(session)
	}
}

// Stops server's net.Listener and calls OnServerStop on all loaded plugins.
func (s *Server) Stop() (err error) {
	// dcPacket := (&Writer{}).SetId(2).
	// 	WriteByte(0).
	// 	WriteString("Server is shutting down.").
	// 	Data()

	// for _, s := range s.sessions {
	// 	if _, err := s.client.Write(dcPacket); err != nil {
	// 		s.client.Close()
	// 	}
	// 	if s.remote != nil {
	// 		s.remote.Close()
	// 	}
	// }

	return s.l.Close()
}

func (s *Server) handleSession(session *session) {
	fromclient, fromserver, cleanup := s.Handler(session)

	var stopped atomic.Bool
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		for {
			if stopped.Load() {
				break
			}

			p, err := ReadPacket(session.client)
			if err != nil {
				break
			}

			if fromclient != nil {
				fromclient(p)
			}
		}

		session.client.Close() //nolint:errcheck
		stopped.Store(true)
		wg.Done()
	}()

	go func() {
		for {
			if stopped.Load() {
				break
			}

			if session.remote == nil {
				continue
			}

			p, err := ReadPacket(session.remote)
			// if err == io.EOF {
			if err != nil {
				break
			}

			if fromserver != nil {
				fromserver(p)
			}
		}

		if session.remote != nil {
			session.remote.Close() //nolint:errcheck
		}

		stopped.Store(true)
		wg.Done()
	}()

	wg.Wait()

	if cleanup != nil {
		cleanup()
	}
}
