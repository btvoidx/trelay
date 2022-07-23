package trelay

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/btvoidx/trelay/packet"
)

type trelay struct {
	opts Options

	listener net.Listener
	plugins  []*lplugin
	players  []TPlayer
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

	return &trelay{
		opts: o,

		players: make([]TPlayer, o.MaxPlayers),
	}
}

// Returns Options used to create this trelay server
func (t *trelay) Options() Options { return t.opts }

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

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	lp := &lplugin{}

	if err := lp.compile(f, path); err != nil {
		return err
	}

	if err := lp.load(t); err != nil {
		return err
	}

	t.plugins = append(t.plugins, lp)
	return nil
}

// Loads a lua plugin just as LoadPlugin, but from fs.FS
func (t *trelay) LoadPluginFS(fs fs.FS, path string) error {
	fileOrDir, err := fs.Open(path)
	if os.IsNotExist(err) {
		return err
	}
	defer fileOrDir.Close()

	info, err := fileOrDir.Stat()
	if err != nil {
		return err
	}

	if info.IsDir() {
		path = filepath.Join(path, "init.lua")
	}

	f, err := fs.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	lp := &lplugin{}

	if err := lp.compile(f, path); err != nil {
		return err
	}

	if err := lp.load(t); err != nil {
		return err
	}

	t.plugins = append(t.plugins, lp)
	return nil
}

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

		go func() {
			if tpl, ok := t.handshake(nc); ok {
				t.connectionLoop(tpl)
			}
		}()
	}

	return nil
}

// Kicks all connected players and stops the proxy
func (t *trelay) Stop() {
	for _, tpl := range t.players {
		if tpl == nil {
			continue
		}
		tpl.Disconnect("Server stopped")
	}

	t.listener.Close()
}

func (t *trelay) getOpenPlayerSlot() int {
	for i, tpl := range t.players {
		if tpl == nil {
			return i
		}
	}

	return -1
}

func (t *trelay) handshake(nc net.Conn) (*tplayer, bool) {
	fmt.Printf("%s is connecting\n", nc.RemoteAddr())

	var tpl *tplayer = &tplayer{
		rconn: nc,
	}

	disconnect := func() { tpl.Disconnect("Possible network tampering detected") }

	for {
		var err error
		p, err := packet.ReadPacket(nc)
		if err != nil {
			fmt.Printf("%s failed handshake\n", nc.RemoteAddr())
			nc.Close()
			return nil, false
		}

		// Stage 1: Client sends a bunch of info: PlayerInfo [4], ClientUUID [5], PlayerHP [16], PlayerMana [42],
		// PlayerBuffs [50], PlayerInventorySlot [5] and finishes with RequestWorldData [6]
		switch p.Type() {
		case packet.ConnectRequest:
			tpl.p_version = p

			var pw packet.Writer
			pw.SetType(packet.SetUserSlot)
			pw.WriteByte(0) // This is player index, which client now will use to send data about itself

			if _, err := nc.Write(pw.Packet().Data()); err != nil {
				disconnect()
				continue
			}

		case packet.PlayerInfo:
			if i, err := p.ReadByte(); err != nil || i != 0 {
				disconnect()
				continue
			}

			tpl.p_info = p

		case packet.ClientUUID:
			tpl.p_uuid = p

		case packet.PlayerHP, packet.PlayerMana, packet.PlayerInventorySlot:
			tpl.p_other = append(tpl.p_other, p)

		case packet.RequestWorldInfo:
			if tpl.Uuid() == "" || tpl.Name() == "" || tpl.Version() == "" {
				disconnect()
				continue
			}

			tpl.index = t.getOpenPlayerSlot()

			if tpl.index < 0 {
				tpl.Disconnect("Server is full")
				return nil, false
			}

			t.players[tpl.index] = tpl
			return tpl, true
		}
	}
}

func (t *trelay) connectionLoop(tpl *tplayer) {
	fmt.Printf("%s (%s) has connected!\n", tpl.rconn.RemoteAddr(), tpl.Name())
	t.callPlugins("on_connect", func(lp *lplugin) int {
		t := lp.LState.NewTable()
		t.RawSetString("player", lp.LState.ToNumber(tpl.index))
		lp.LState.Push(t)
		return 1
	})

	defer fmt.Printf("%s (%s) has disconnected!\n", tpl.rconn.RemoteAddr(), tpl.Name())
	defer t.callPlugins("on_disconnect", func(lp *lplugin) int {
		t := lp.LState.NewTable()
		t.RawSetString("player", lp.LState.ToNumber(tpl.index))
		lp.LState.Push(t)
		return 1
	})

	// todo
	if err := tpl.ChangeServer("localhost:7778"); err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		time.Sleep(time.Second * 5)
		println("CONNECTING")
		err := tpl.ChangeServer("localhost:7779")
		if err != nil {
			fmt.Println(err)
		}
	}()

	var shouldClose bool
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		defer func() {
			shouldClose = true
			tpl.rconn.Close()
		}()

		for {
			if shouldClose {
				break
			}

			if tpl.sconn == nil {
				continue
			}

			p, err := packet.ReadPacket(tpl.rconn)
			if err == io.EOF && tpl.sconn == nil {
				break
			}

			if err != nil {
				continue
			}

			if p.Type() == packet.PlayerActive {
				v, _ := p.ReadByte()
				tpl.knownPlayers[v] = struct{}{}
				println("New player is known!")
			} else if p.Type() == packet.NPCUpdate {
				v, _ := p.ReadInt16()
				tpl.knownNPCs[v] = struct{}{}
				println("New npc is known!")
			} else if p.Type() == packet.UpdateItemDrop {
				v, _ := p.ReadInt16()
				tpl.knownItems[v] = struct{}{}
				println("New item is known!")
			}

			if err != nil {
				tpl.Disconnect("Network tampering detected")
				break
			}

			// var handled bool

			// t.callPlugins("on_player_packet", func(lp *lplugin) int {
			// 	t := lp.LState.NewTable()
			// 	t.RawSetString("player", tpl.toTable(lp.LState)) // todo: pull from plugin's cache
			// 	t.RawSetString("handled", luafnSetPacketHandled(lp.LState, &handled))
			// 	lp.LState.Push(t)
			// 	return 1
			// })

			// if handled {
			// 	continue
			// }

			if _, err := tpl.sconn.Write(p.Data()); err != nil {
				break
			}
		}
	}()

	go func() {
		defer wg.Done()
		defer func() {
			shouldClose = true
			if tpl.sconn != nil {
				tpl.sconn.Close()
			}
		}()

		for {
			if shouldClose {
				break
			}

			if tpl.sconn == nil {
				continue
			}

			p, err := packet.ReadPacket(tpl.sconn)
			if err == io.EOF && tpl.sconn != nil {
				break
			}

			if err != nil {
				continue
			}

			if _, err := tpl.rconn.Write(p.Data()); err != nil {
				break
			}
		}
	}()

	wg.Wait()
}

// Calls all plguins asynchronously and waits for all of them
func (t *trelay) callPlugins(fn string, ctx func(*lplugin) int) {
	var wg sync.WaitGroup

	for _, L := range t.plugins {
		wg.Add(1)
		go func(L *lplugin) {
			L.Call(fn, ctx)
			wg.Done()
		}(L)
	}

	wg.Wait()
}
