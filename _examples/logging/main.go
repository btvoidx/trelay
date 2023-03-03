package main

import (
	"fmt"

	"github.com/btvoidx/trelay"
)

func main() {
	err := trelay.ListenAndServe(":7777", &LoggingHandler{&trelay.Direct{Addr: "213.108.4.58:7777"}})
	if err != nil {
		fmt.Printf("An error occured when starting the server: %v\n", err)
	}
}

type LoggingHandler struct{ trelay.Handler }

func (h *LoggingHandler) ClientPacket(s trelay.Session, p trelay.Packet) bool {
	fmt.Printf("sent id:%-3d len:%d\n", p.Id(), p.Length())
	return h.Handler.ClientPacket(s, p)
}

func (h *LoggingHandler) RemotePacket(s trelay.Session, p trelay.Packet) bool {
	fmt.Printf("rcvd id:%-3d len:%d\n", p.Id(), p.Length())
	return h.Handler.RemotePacket(s, p)
}
