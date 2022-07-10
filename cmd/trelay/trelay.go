package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/btvoidx/trelay"
	"github.com/spf13/pflag"
)

func main() {
	var (
		plugins []string
		addr    string
	)

	pflag.StringArrayVarP(&plugins, "plugin", "p", []string{}, "a path to plugin directory; can be given multiple times to load multiple plugins")
	pflag.StringVarP(&addr, "address", "a", "0.0.0.0:7777", "network address to start trelay server on")
	pflag.Parse()

	trelay := trelay.NewTrelayServer(trelay.Options{
		MaxPlayers: 1984,
	})

	for _, path := range plugins {
		fmt.Printf("Loading plugin: %s\n", path)
		if err := trelay.LoadPlugin(path); err != nil {
			fmt.Printf("An error occured while loading plugin '%s': %v\n", path, err)
		}
	}

	go func() {
		fmt.Printf("Starting trelay on %s\n", addr)

		if err := trelay.Start(addr); err != nil {
			fmt.Printf("Alas, an error: %v\n", err)
		}
	}()

	go func() {
		// todo fmt.Printf("Starting rest api on %s\n", addr)
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	trelay.Stop()
}
