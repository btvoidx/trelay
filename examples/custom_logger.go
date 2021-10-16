package main

import (
	"os"
	"os/signal"

	"github.com/btvoidx/trelay"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Custom logger
	log := log.New().WithField("Custom Field", true)

	server := trelay.NewServer("localhost:7777", "terraria.tk:7777").
		SetLogger(log) // Use custom logger

	err := server.Start()
	if err != nil {
		log.Fatal("An error occured when starting the server: ", err.Error())
	}

	defer server.Stop()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
