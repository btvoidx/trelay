package main

import (
	"log"

	"github.com/btvoidx/trelay"
)

func main() {
	if err := trelay.ListenAndServe(":7777", trelay.Direct("213.108.4.58:7777", 1)); err != nil {
		log.Fatalf("An error occured when starting the server: %v", err)
	}
}
