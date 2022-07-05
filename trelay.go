package trelay

import (
	"errors"
	"fmt"
	"math"
	"net"
	"os"
	"path/filepath"
	"sync"

	"github.com/btvoidx/trelay/packet"
	lua "github.com/yuin/gopher-lua"
)

const globalLuaTableName = "trelay"

type trelay struct {
	opts Options

	listener net.Listener
	plugins  []*luaplugin
	players  []TPlayer
	servers  []TServer
}

type Options struct {
	// Maximum number of players trelay allows to connect.
	// Default: 256
	//
	// Negative values will prevent server from starting.
	MaxPlayers int
}

func NewTrelayServer(opts ...Options) *trelay {
	if len(opts) == 0 {
		opts = append(opts, Options{})
	}

	o := opts[0]

	if o.MaxPlayers == 0 {
		o.MaxPlayers = math.MaxUint8
	}

	return &trelay{opts: o}
}

// Loads a lua plugin at the given path.
// If it's a file, it gets loaded.
// If it's a directory, "init.lua" from it gets loaded.
func (t *trelay) LoadPlugin(path string) error {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return err
	}

	if info.IsDir() {
		path = filepath.Join(path, "init.lua")
	}

	L := lua.NewState()
	L.G.Global.RawSetString(globalLuaTableName, L.NewTable())

	if err := L.DoFile(path); err != nil {
		return err
	}

	t.plugins = append(t.plugins, &luaplugin{
		LState: L,
		Mutex:  sync.Mutex{},
	})
	return nil
}

// Same as LoadPlugin, but loads from fs.FS
// func (t *trelay) LoadPluginFS(fs fs.FS, path string) error {
// 	f, err := fs.Open(path)
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()

// 	info, err := f.Stat()
// 	if err != nil {
// 		return err
// 	}

// 	if info.IsDir() {
// 		path = filepath.Join(path, "init.lua")
// 	}

// 	source := make([]byte, 0, info.Size())
// 	if n, err := f.Read(source); err != nil || n < int(info.Size()) {
// 		return err
// 	}

// 	L := lua.NewState()
// 	if err := L.DoString(string(source)); err != nil {
// 		return err
// 	}

// 	t.plugins = append(t.plugins, L)
// 	return nil
// }

func (t *trelay) Start(addr string) error {
	if t.opts.MaxPlayers < 0 {
		return fmt.Errorf("options.MaxPlayers is %d; must be positive", t.opts.MaxPlayers)
	}

	l, err := net.Listen("tcp4", addr)
	if err != nil {
		return err
	}

	t.listener = l

	for {
		nc, err := t.listener.Accept()
		if err != nil && errors.Is(err, net.ErrClosed) {
			break
		}

		if t.handshake(nc) {
			t.connectionLoop(nc)
		}
	}

	return nil
}

// Kicks all connected players and stops the proxy
func (t *trelay) Stop() {
	for _, tpl := range t.players {
		tpl.Disconnect("Server stopped")
	}
}

func (t *trelay) getOpenPlayerSlot() int {
	for i, tpl := range t.players {
		if tpl == nil {
			return i
		}
	}

	t.players = append(t.players, nil)
	return len(t.players)
}

func (t *trelay) handshake(conn net.Conn) bool {
	fmt.Printf("%s is connecting\n", conn.RemoteAddr())

	var stage int = 0
	var tpl *tplayer = &tplayer{
		conn: conn,
	}

	disconnect := func() { tpl.Disconnect("Possible network tampering detected") }

	for {
		var err error
		p, err := packet.ReadPacket(conn)
		if err != nil {
			fmt.Printf("%s disconnected\n", conn.RemoteAddr())
			conn.Close()
			return false
		}

		switch stage {
		// Stage 0: Client sends ConnectRequest [1], waits for SetUserSlot [3]
		case 0:
			if p.Type() != packet.ConnectRequest {
				disconnect()
				continue
			}

			tpl.version, err = p.ReadString()
			if err != nil {
				disconnect()
				continue
			}

			var pw packet.Writer
			pw.SetType(packet.SetUserSlot)
			pw.WriteByte(0) // This is player index, which client now will use to send data about itself

			if _, err := conn.Write(pw.Packet().Data()); err != nil {
				disconnect()
				continue
			}

			stage++

		// Stage 1: Client sends a bunch of info: PlayerInfo [4], ClientUUID [5], PlayerHP [16], PlayerMana [42],
		// PlayerBuffs [50], PlayerInventorySlot [5] and finishes with RequestWorldData [6]
		case 1:
			switch p.Type() {
			case packet.PlayerInfo:
				if i, err := p.ReadByte(); err != nil || i != 0 {
					disconnect()
					continue
				}

				if _, err = p.ReadBytes(2); err != nil {
					disconnect()
					continue
				}

				if tpl.name, err = p.ReadString(); err != nil {
					disconnect()
					continue
				}

			case packet.ClientUUID:
				if tpl.uuid, err = p.ReadString(); err != nil {
					disconnect()
					continue
				}

			case packet.RequestWorldData:
				if tpl.uuid == "" || tpl.name == "" || tpl.version == "" {
					disconnect()
					continue
				}

				fmt.Printf("A player has requested world data! %v\n", tpl)
				disconnect()
				continue
			}

		case 2:
			// todo t.players = append(t.players, tpl)
		}

	}
}

func (t *trelay) connectionLoop(conn net.Conn) {
	// go func() {
	// 	for {
	// 		p, err := packet.ReadPacket(conn)
	// 		// if err == io.EOF {
	// 		// 	return
	// 		// }
	// 		if err != nil {
	// 			continue
	// 		}

	// 		ctx := &PacketContext{
	// 			packet: p,
	// 		}

	// 		if ctx.handled {
	// 			continue
	// 		}

	// 		// if session.Server.Write(p) != nil {
	// 		// 	break
	// 		// }
	// 	}
	// }()

	// go func() {
	// 	for {
	// 		p, err := session.Server.Read()
	// 		if err != nil {
	// 			break
	// 		}

	// 		ctx := &PacketContext{
	// 			packet:  p,
	// 			session: session,
	// 		}

	// 		for _, plugin := range s.plugins {
	// 			plugin.OnServerPacket(ctx)
	// 		}

	// 		if ctx.handled {
	// 			continue
	// 		}

	// 		if session.Client.Write(p) != nil {
	// 			break
	// 		}
	// 	}
	// }()
}

func (t *trelay) callPlugins(fname string) {
	for _, L := range t.plugins {
		L.Lock()
		defer L.Unlock()

		gt, ok := L.G.Global.RawGetString(globalLuaTableName).(*lua.LTable)
		if !ok {
			println("ERROR TRELAY NOT A TABLE: ", fname) //todo logging
			return
		}

		fn, ok := gt.RawGetString(fname).(*lua.LFunction)
		if !ok {
			return
		}

		L.Push(fn)
		if err := L.PCall(0, lua.MultRet, nil); err != nil {
			println("ERROR IN A FUNCTION: ", fname) //todo logging
			return
		}
	}
}
