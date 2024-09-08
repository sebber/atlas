package main

import (
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net"
	"sync"
	"time"

	ping "github.com/sebber/atlas/internal"
)

type Client struct {
	Nickname string
	Conn     net.Conn
}

var clients = make(map[string]Client)
var clientsMutex sync.Mutex

func main() {
	slog.Info("Atlas starting")

	port := flag.Int("port", 8123, "Port number for the atlast server")
	flag.Parse()

	address := fmt.Sprintf(":%d", *port)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		slog.Error("Failed to listen", slog.Int("port", *port), slog.Any("error", err))
		return
	}
	defer listener.Close()

	slog.Info("Server is running", slog.Int("port", *port))

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

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		slog.Error("Error reading Nickname", slog.Any("error", err))
		return
	}
	clientName := string(buf[:n])

	clientsMutex.Lock()
	clients[clientName] = Client{Nickname: clientName, Conn: conn}
	slog.Info("Connected as", slog.String("Nickname", clientName))
	logActiveClients()
	clientsMutex.Unlock()

	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				clientsMutex.Lock()
				delete(clients, clientName)
				slog.Info("Client disconnected")
				logActiveClients()
				clientsMutex.Unlock()
				break
			}

			slog.Error("Error reading from connection", slog.Any("error", err))
			return
		}

		msg, err := ping.Deserialize(buf[:n])
		if err != nil {
			slog.Error("Failed to deserialize ping", slog.Any("error", err))
			return
		}

		slog.Info("Received Ping", slog.Int64("timestamp", msg.Timestamp))

		err = sendPong(conn)
		if err != nil {
			slog.Error("Failed to send Pong", slog.Any("error", err))
		}
	}
}

func sendPong(conn net.Conn) error {
	timestamp := time.Now().Unix()
	pongMessage, err := ping.Serialize(timestamp)
	if err != nil {
		return err
	}

	_, err = conn.Write(pongMessage)
	return err
}

func logActiveClients() {
	slog.Info("Active clients:")
	for name := range clients {
		slog.Info("client", slog.String("Nickname", name))
	}
}
