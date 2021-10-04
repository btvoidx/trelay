package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/btvoidx/trelay"
	log "github.com/sirupsen/logrus"
)

func main() {
	server := trelay.NewServer("0.0.0.0:7777", "terraria.tk:7777")
	server.SetLogger(
		log.New().WithField("Custom Field", true),
	)

	err := server.Start()
	if err != nil {
		fmt.Println("An error occured when starting the server: ", err.Error())
		os.Exit(1)
	}

	defer server.Stop()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
