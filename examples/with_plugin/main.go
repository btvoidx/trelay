package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/btvoidx/trelay"
)

type plugin struct{}

func (*plugin) Name() string { return "ExamplePlugin" }

func (p *plugin) OnLoad(s trelay.Server) trelay.Plugin {
	fmt.Println("plugin.OnLoad")
	fmt.Printf("plugin loaded by server:%d\n", s.Id())
	return p
}

func (*plugin) OnServerStart(s trelay.Server)   { fmt.Println("plugin.OnServerStart") }
func (*plugin) OnServerStop(s trelay.Server)    { fmt.Println("plugin.OnServerStop") }
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

func main() {
	server := trelay.NewServer("localhost:7777", "terraria.tk:7777").
		LoadPlugin(&plugin{})

	err := server.Start()
	if err != nil {
		log.Fatalf("An error occured when starting the server: %s", err.Error())
	}

	defer server.Stop()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
