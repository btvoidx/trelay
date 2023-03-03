package main

import (
	"log"

	"github.com/btvoidx/trelay"
)

func main() {
	err := trelay.ListenAndServe(":7777", &trelay.Direct{Addr: "213.108.4.58:7777", MaxPlayers: 1})
	if err != nil {
		log.Fatalf("An error occured when starting the server: %v", err)
	}
}
