package client

import (
	"log/slog"
	"net"
	"time"

	messaging "github.com/sebber/atlas/internal/messaging"
)

type Client struct {
	Conn   net.Conn
	ConnId string
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Connect(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		slog.Error("Failed to dial", slog.String("address", addr), slog.Any("error", err))
		return err
	}
	defer conn.Close()

	msg, err := messaging.ReceiveMessage(conn)
	if err != nil {
		slog.Error("Couldn't read message", slog.Any("error", err))
		return err
	}

	switch m := msg.(type) {
	case *messaging.ConnStartMessage:
		slog.Info("got a ConnStartMessage", slog.Any("type", m.MessageType()), slog.String("id", m.Id))
	default:
		slog.Info("Got unknown message")
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:

			pingMsg := messaging.PingMessage{Timestamp: time.Now().Unix()}

			err := messaging.SendMessage(conn, &pingMsg)
			if err != nil {
				slog.Error("Failed to send Ping", slog.Any("error", err))
				return err
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
