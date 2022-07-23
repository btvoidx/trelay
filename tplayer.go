package trelay

import (
	"fmt"
	"net"

	"github.com/btvoidx/trelay/packet"
)

// TPlayer represents both a client and server connections, as well as their actual Terraria player
type TPlayer interface {
	// This player's index
	Index() int

	// Player game version as sent by Connect Request [1] packet
	Version() string

	// Client UUID as sent by Client UUID [68] packet
	Uuid() string

	// Player name, as sent by Player Info [4]
	Name() string

	// The server the player is currently connected to
	//
	// Use ChangeServer to connect player to another server
	// Server() *TServer

	// Abruptly disconnects the player from terraria server and from trelay.
	// Default reason is "Disconnected", only first reason is taken. Disconnect reason is not given to the server.
	Disconnect(reason string, a ...any)
}

type tplayer struct {
	index int
	// remote connection and server connection
	rconn, sconn net.Conn

	p_version *packet.Packet
	p_uuid    *packet.Packet
	p_info    *packet.Packet
	p_other   []*packet.Packet

	knownPlayers map[byte]struct{}
	knownNPCs    map[int16]struct{}
	knownItems   map[int16]struct{}
}

func (tpl *tplayer) Index() int { return tpl.index }

func (tpl *tplayer) Version() string {
	tpl.p_version.ResetHead()
	if v, err := tpl.p_version.ReadString(); err == nil {
		return v
	}

	return ""
}

func (tpl *tplayer) Uuid() string {
	tpl.p_uuid.ResetHead()
	if v, err := tpl.p_uuid.ReadString(); err == nil {
		return v
	}

	return ""
}

func (tpl *tplayer) Name() string {
	tpl.p_info.ResetHead()
	if err := tpl.p_info.AdvanceHead(3); err != nil {
		return ""
	}

	if v, err := tpl.p_info.ReadString(); err == nil {
		return v
	}

	return ""
}

// // The server the player is currently connected to
// //
// // Use ChangeServer to connect player to another server
// func (tp *tplayer) Server() *TServer { return p.server }

// // Options struct to Player.ChangeServer()
// type ChangeServerOptions struct {
// 	Password string
// }

// Connect player to other server. Connection with current server is not closed until
// new connection is fully established. Both synchronous and asynchronous use is ok? Only synchronous tested.
//
// Use options to connect to a password-protected server
func (tp *tplayer) ChangeServer(addr string) error {
	nc, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	if _, err := nc.Write(tp.p_version.Data()); err != nil {
		nc.Close()
		return err
	}

	prevServer := tp.sconn
	tp.sconn = nil

	// // Send all player info
	// {
	// 	var userslot byte
	// 	if p, err := packet.ReadPacket(nc); err != nil {
	// 		nc.Close()
	// 		return err
	// 	} else if p.Type() != packet.SetUserSlot {
	// 		nc.Close()
	// 		return err
	// 	} else {
	// 		if userslot, err = p.ReadByte(); err != nil {
	// 			return err
	// 		}
	// 		tp.rconn.Write(p.Data()) //nolint:errcheck
	// 	}

	// 	// Send player info packet to server replacing first byte with user slot
	// 	// the "{" is here to scope out data
	// 	{
	// 		data := tp.p_info.Data()
	// 		data[4] = userslot // 4 is where PlayerID is
	// 		if _, err := nc.Write(data); err != nil {
	// 			nc.Close()
	// 			return err
	// 		}
	// 	}

	// 	// Send all other packets
	// 	for _, p := range tp.p_other {
	// 		data := p.Data()
	// 		data[4] = userslot // same as above
	// 		if _, err := nc.Write(data); err != nil {
	// 			nc.Close()
	// 			return err
	// 		}
	// 	}
	// }

	// Clear known players, NPCs, items
	{
		for v := range tp.knownPlayers {
			data := (&packet.Writer{}).
				SetType(packet.PlayerActive).
				WriteByte(v).
				WriteByte(0).
				Packet().
				Data()

			tp.rconn.Write(data) //nolint:errcheck
			tp.knownPlayers = make(map[byte]struct{})
		}

		for v := range tp.knownNPCs {
			data := (&packet.Writer{}).
				SetType(packet.NPCUpdate).
				WriteInt16(v).
				// Enable monospace font to properly see this comment:
				//                Position X; Position Y; Velocity X; Velocity Y; Target    ; Flag; NPC Net ID; Life
				WriteBytes([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0}).
				Packet().
				Data()

			tp.rconn.Write(data) //nolint:errcheck
			tp.knownNPCs = make(map[int16]struct{})
		}

		for v := range tp.knownItems {
			data := (&packet.Writer{}).
				SetType(packet.UpdateItemDrop).
				WriteInt16(v).
				// Enable monospace font to properly see this comment:
				//                Position X; Position Y; Velocity X; Velocity Y; Stck; P;
				WriteBytes([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}).
				Packet().
				Data()

			tp.rconn.Write(data) //nolint:errcheck
			tp.knownItems = make(map[int16]struct{})
		}
	}

	if prevServer != nil {
		prevServer.Close()
	}

	tp.sconn = nc
	return nil

	// Request world info and get spawn section
	{
		// A literal for (&packet.Writer{}).SetType(packet.RequestWorldInfo).Packet().Data()
		requestWorldInfo := []byte{3, 0, byte(packet.RequestWorldInfo)}

		if _, err := nc.Write(requestWorldInfo); err != nil {
			nc.Close()
			return err
		}

		if p, err := packet.ReadPacket(nc); err != nil {
			nc.Close()
			return err
		} else if p.Type() != packet.WorldInfo {
			nc.Close()
			return err
		} else {
			tp.rconn.Write(p.Data()) //nolint:errcheck
		}

		// todo: replace with packet literal
		// if _, err := nc.Write((&packet.Writer{}).SetType(packet.GetSectionOrRequestSync).WriteInt16(-1).WriteInt16(-1).Packet().Data()); err != nil {
		// 	nc.Close()
		// 	return err
		// }

		tp.sconn = nc

		// var spawnX, spawnY int16
		// if p, err := packet.ReadPacket(nc); err != nil {
		// 	nc.Close()
		// 	return err

		// } else if p.Type() != packet.WorldInfo {
		// 	nc.Close()
		// 	return err

		// } else {
		// 	p.AdvanceHead(4 + 1 + 1 + 2 + 2) // time, mooninfo, moonphase, max tiles x, max tiles y
		// 	if spawnX, err = p.ReadInt16(); err != nil {
		// 		nc.Close()
		// 		return err
		// 	}
		// 	if spawnY, err = p.ReadInt16(); err != nil {
		// 		nc.Close()
		// 		return err
		// 	}
		// }
	}

	return nil
}

func (tpl *tplayer) Disconnect(reason string, a ...any) {
	reason = fmt.Sprintf(reason, a...)

	var pw packet.Writer
	pw.SetType(packet.Disconnect)
	pw.WriteByte(0)
	pw.WriteString(reason)

	tpl.rconn.Write(pw.Packet().Data()) //nolint:errcheck
	tpl.rconn.Close()
}
