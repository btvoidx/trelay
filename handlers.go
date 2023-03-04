package trelay

import (
	"fmt"
	"net"
	"sync"
)

var _ Handler = (*Direct)(nil)

// The simplest Handler. Connects the client to remote server at specified address and makes sure
// total player count is not exceeded.
// Default disconnect messages can be changed by changing ConnectionFailed and ServerIsFull global variables.
type Direct struct {
	Addr string

	// Zero means no limit, set to negative to allow nobody
	MaxPlayers int64
	// Disconnect message when connection to Addr fails
	ConnectionFailed string
	// Disconnect message when server is full
	ServerIsFull string

	currentPlayers int64
	mu             sync.Mutex
}

func (h *Direct) ClientConnect(s Session) {
	h.mu.Lock()

	max := h.MaxPlayers // copy value because it can be changed from another goroutine
	if max != 0 && h.currentPlayers >= max {
		h.mu.Unlock()

		msg := h.ServerIsFull
		if msg == "" {
			msg = "trelay: server is full"
		}

		s.Client().WritePacket(new(Builder).
			SetId(2).
			WriteByte(0).
			WriteString(msg))
		s.Close()
		return
	}

	h.currentPlayers += 1
	h.mu.Unlock()

	remote, err := net.Dial("tcp", h.Addr)
	if err != nil {
		msg := h.ConnectionFailed
		if msg == "" {
			msg = fmt.Sprintf("trelay: failed to connect to remote server at %s", h.Addr)
		}

		s.Client().WritePacket(new(Builder).
			SetId(2).
			WriteByte(0).
			WriteString(msg))
		s.Close()
		return
	}

	s.SetRemote(remote)
}

func (h *Direct) ClientDisconnect(s Session) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.currentPlayers -= 1
}

func (h *Direct) ClientPacket(Session, Packet) bool { return false }
func (h *Direct) RemotePacket(Session, Packet) bool { return false }
