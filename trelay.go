package trelay

import (
	"net"
	"sync/atomic"
)

func ListenAndServe(addr string, h SessionHandler) error {
	return (&Server{Addr: addr, Handler: h}).ListenAndServe()
}

// A simplest SessionHandler. Connects the client to the remote server at specified address and makes sure
// total player count is not exceeded.
func Direct(addr string, maxPlayers int64) SessionHandler {
	dcPacket := (&Writer{}).SetId(2).
		WriteByte(0).
		WriteStringf("trelay: couldn't connect to %s", addr).
		Data()

	maxPlayersPacket := (&Writer{}).SetId(2).
		WriteByte(0).
		WriteString("trelay: server is full").
		Data()

	currentPlayers := atomic.Int64{}

	return func(s Session) (fromclient func(Packet), fromserver func(Packet), cleanup func()) {
		if currentPlayers.Load() >= maxPlayers {
			s.Client().Write(maxPlayersPacket) //nolint:errcheck
			s.Client().Close()
			return
		}

		currentPlayers.Add(1)
		cleanup = func() { currentPlayers.Add(-1) }

		if s.Remote() == nil {
			serv, err := net.Dial("tcp", addr)
			if err != nil {
				s.Client().Write(dcPacket) //nolint:errcheck
				s.Client().Close()
				return
			}

			s.SetRemote(serv)
		}

		fromclient = func(p Packet) { s.Remote().Write(p.Buffer()) } //nolint:errcheck
		fromserver = func(p Packet) { s.Client().Write(p.Buffer()) } //nolint:errcheck

		return
	}
}
