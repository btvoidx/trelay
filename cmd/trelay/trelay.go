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

	pflag.StringArrayVarP(&plugins, "plugins", "p", []string{"plugins"}, "list of directories to pull plugins from or paths directly to plugins")
	pflag.StringVarP(&addr, "address", "a", "0.0.0.0:7777", "network address to start trelay server on")
	pflag.Parse()

	server := trelay.NewTrelayServer(trelay.Options{
		MaxPlayers: 1984,
	})

	for _, path := range plugins {
		fmt.Printf("Loading plugin: %s\n", path)
		if err := server.LoadPlugin(path); err != nil {
			fmt.Printf("An error occured while loading plugin '%s': %v\n", path, err)
		}
	}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		os.Exit(1)
	}()

	fmt.Printf("Starting trelay on %s\n", addr)
	if err := server.Start(addr); err != nil {
		fmt.Printf("Alas, an error: %v\n", err)
	}
}
