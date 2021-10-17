package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/btvoidx/trelay"
)

func main() {
	server := trelay.NewServer("localhost:7777", "terraria.tk:7777").
		LoadPlugin(GetExamplePlugin) // look plugin.go

	err := server.Start()
	if err != nil {
		log.Fatalf("An error occured when starting the server: %s", err.Error())
	}

	defer server.Stop()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
