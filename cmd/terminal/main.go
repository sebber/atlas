package main

import (
	"flag"
	"fmt"
	"log/slog"

	"github.com/sebber/atlas/internal/client"
)

func main() {
	slog.Info("Terminal starting")

	port := flag.Int("port", 8123, "Port number for the atlast server")
	flag.Parse()

	address := fmt.Sprintf("localhost:%d", *port)

	client := client.NewClient()

	client.Connect(address)
}
