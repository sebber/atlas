package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net"
	"time"

	messaging "github.com/sebber/atlas/internal/messaging"
)

func main() {
	slog.Info("Terminal starting")

	port := flag.Int("port", 8123, "Port number for the atlast server")
	nickname := flag.String("nickname", "unnamed", "Nickname for the whatevs")
	flag.Parse()

	if *nickname == "unnamed" {
		slog.Error("Nickname required")
		return
	}
	slog.Info("Connecting as", slog.String("nickname", *nickname))

	address := fmt.Sprintf("localhost:%d", *port)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		slog.Error("Failed to dial", slog.String("address", address), slog.Any("error", err))
		return
	}
	defer conn.Close()

	msg, err := messaging.ReceiveMessage(conn)
	if err != nil {
		slog.Error("Couldn't read message", slog.Any("error", err))
		return
	}

	switch m := msg.(type) {
	case *messaging.ConnStartMessage:
		slog.Info("got a ConnStartMessage", slog.Any("type", m.MessageType()), slog.String("id", m.Id))
	default:
		slog.Info("Got unknown message")
	}

	// _, err = conn.Write([]byte(*nickname))
	// if err != nil {
	// 	slog.Error("Failed identifying with nickname", slog.Any("error", err))
	// 	return
	// }

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:

			pingMsg := messaging.PingMessage{Timestamp: time.Now().Unix()}

			err := messaging.SendMessage(conn, &pingMsg)
			if err != nil {
				slog.Error("Failed to send Ping", slog.Any("error", err))
				return
			}

			// buf := make([]byte, 1024)
			// n, err := conn.Read(buf)
			// if err != nil {
			// 	slog.Error("Failed to read Pong", slog.Any("error", err))
			// 	return
			// }
			//
			// pongMsg, err := messaging.Deserialize((buf[:n]))
			// if err != nil {
			// 	slog.Error("Failed to deserialize pong", slog.Any("error", err))
			// 	return
			// }
			//
			// slog.Info("Received Pong", slog.Int64("timestamp", pongMsg.Timestamp))
		}
	}
}
