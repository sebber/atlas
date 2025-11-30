package main

import (
	"log/slog"

	server "github.com/sebber/atlas/internal/server"
)

func main() {
	slog.Info("Atlas starting")

	srv := server.NewServer()
	srv.Start()
}
