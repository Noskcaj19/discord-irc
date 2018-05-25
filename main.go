package main

import (
	"fmt"
	"os"
)

func main() {
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		fmt.Fprintln(os.Stderr, "`DISCORD_TOKEN` not set, exiting")
		os.Exit(1)
	}

	server := NewServer("nos")
	server.Run()
}
