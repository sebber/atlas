package main

import (
	"flag"
	"log/slog"

	server "github.com/sebber/atlas/internal/server"
)

func main() {
	slog.Info("Atlas starting")

	port := flag.Int("port", 8123, "Port number for the atlast server")
	flag.Parse()

	srv := server.NewServer(*port)
	srv.Start()
}
