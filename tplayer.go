package trelay

import (
	"fmt"
	"net"

	"github.com/btvoidx/trelay/packet"
	lua "github.com/yuin/gopher-lua"
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
	conn   net.Conn
	server *TServer

	index   int
	version string
	uuid    string
	name    string
}

func (tpl *tplayer) Index() int      { return tpl.index }
func (tpl *tplayer) Version() string { return tpl.version }
func (tpl *tplayer) Uuid() string    { return tpl.uuid }
func (tpl *tplayer) Name() string    { return tpl.name }

// // The server the player is currently connected to
// //
// // Use ChangeServer to connect player to another server
// func (tp *tplayer) Server() *TServer { return p.server }

// // Options struct to Player.ChangeServer()
// type ChangeServerOptions struct {
// 	Password string
// }

// // Connect player to other server. Connection with current server is not closed until
// // new connection is fully established. Both synchronous and asynchronous use is ok.
// //
// // Use options to connect to a password-protected server
// func (tp *tplayer) ChangeServer(addr string, o ...ChangeServerOptions) {}

func (tpl *tplayer) Disconnect(reason string, a ...any) {
	reason = fmt.Sprintf(reason, a...)

	var pw packet.Writer
	pw.SetType(packet.Disconnect)
	pw.WriteByte(0)
	pw.WriteString(reason)

	tpl.conn.Write(pw.Packet().Data()) //nolint:errcheck
	tpl.conn.Close()
}

func (tpl *tplayer) toTable(L *lua.LState) *lua.LTable {
	t := L.NewTable()

	t.RawSetString("index", lua.LNumber(tpl.Index()))
	t.RawSetString("version", lua.LString(tpl.Version()))
	t.RawSetString("name", lua.LString(tpl.Name()))
	t.RawSetString("uuid", lua.LString(tpl.Uuid()))

	t.RawSetString("disconnect", L.NewFunction(func(l *lua.LState) int {
		tpl.Disconnect(l.CheckString(1))
		return 0
	}))

	return t
}
