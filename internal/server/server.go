package server

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/sebber/atlas/internal/messaging"
)

type Client struct {
	Conn net.Conn
}

var clients = make(map[string]Client)
var clientsMutex sync.Mutex

type Server struct {
	Port int
}

func NewServer(port int) *Server {
	return &Server{Port: port}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
	if err != nil {
		slog.Error("Could not start server", slog.Any("error", err))
		return err
	}
	defer listener.Close()
	slog.Info("Server started")

	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.Error("Failed to accept connection", slog.Any("error", err))
			continue
		}
		slog.Info("New connection established")

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	connId := generateConnectionID()

	connMsg := messaging.ConnStartMessage{Id: connId}
	err := messaging.SendMessage(conn, &connMsg)
	if err != nil {
		slog.Error("Failed to send message", slog.Any("error", err))
	}

	for {
		ctx := context.Background()
		context.WithValue(ctx, "connId", connId)

		msg, err := messaging.ReceiveMessage(conn)
		if err != nil {
			if err == io.EOF {
				break
			}
			slog.Error("Could not read message", slog.Any("error", err))
			continue
		}

		switch msg.(type) {
		case *messaging.PingMessage:
			slog.Info("Got a ping message from", slog.Any("connId", connId))
		default:
			slog.Info("Got unknown message from", slog.Any("connId", connId))
		}
	}
}

func disconnectClient(conn net.Conn, nickname string) {
	clientsMutex.Lock()
	delete(clients, nickname)
	slog.Info("Client disconnected")
	logActiveClients()
	clientsMutex.Unlock()
}

func logActiveClients() {
	slog.Info("Active clients:")
	for name := range clients {
		slog.Info("client", slog.String("Nickname", name))
	}
}

func generateConnectionID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
