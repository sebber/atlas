package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net"
	"time"

	ping "github.com/sebber/atlas/internal"
)

func main() {
	slog.Info("Terminal starting")

	port := flag.Int("port", 8123, "Port number for the atlast server")
	flag.Parse()

	address := fmt.Sprintf("localhost:%d", *port)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		slog.Error("Failed to dial", slog.String("address", address), slog.Any("error", err))
		return
	}
	defer conn.Close()

	timestamp := time.Now().Unix()
	pingMessage, err := ping.Serialize(timestamp)
	if err != nil {
		slog.Error("Failed to serialize timestamp", slog.Int64("timestamp", timestamp), slog.Any("error", err))
		return
	}

	_, err = conn.Write(pingMessage)
	if err != nil {
		slog.Error("Failed to send ping", slog.Any("error", err))
	}
}
