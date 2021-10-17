package main

import (
	"fmt"

	"github.com/btvoidx/trelay"
)

func GetExamplePlugin(s trelay.Server) trelay.Plugin {
	return &plugin{s}
}

type plugin struct{ s trelay.Server }

func (*plugin) Name() string                    { return "ExamplePlugin" }
func (*plugin) OnServerStart()                  { fmt.Println("plugin.OnServerStart") }
func (*plugin) OnServerStop()                   { fmt.Println("plugin.OnServerStop") }
func (*plugin) OnSessionOpen(s trelay.Session)  { fmt.Println("plugin.OnSessionOpen") }
func (*plugin) OnSessionClose(s trelay.Session) { fmt.Println("plugin.OnSessionClose") }

func (*plugin) OnClientPacket(pid trelay.PacketType, packet *trelay.Packet, session trelay.Session) (handled bool) {
	fmt.Println("plugin.OnClientPacket")
	return
}

func (*plugin) OnServerPacket(pid trelay.PacketType, packet *trelay.Packet, session trelay.Session) (handled bool) {
	fmt.Println("plugin.OnServerPacket")
	return
}
