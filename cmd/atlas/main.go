package main

import (
	"context"
	"flag"
	"fmt"
	"io"

	// "io"
	"log/slog"
	"net"
	"sync"
	"time"

	messaging "github.com/sebber/atlas/internal/messaging"
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
	// msg, err := messaging.ReadMessage(conn)
	// if err != nil {
	// if err == io.EOF {
	// disconnectClient(conn, connId)
	// break
	// }

	// slog.Error("Could not read message", slog.Any("error", err))
	// continue
	// }
	//  if err != nil {
	// if err == io.EOF {
	// 	clientsMutex.Lock()
	// 	delete(clients, clientName)
	// 	slog.Info("Client disconnected")
	// 	logActiveClients()
	// 	clientsMutex.Unlock()
	// 	break
	// }
	//
	// slog.Error("Error reading from connection", slog.Any("error", err))
	// return
	// }

	// switch msg.Type {
	// case messaging.MsgTypeAuth:
	// 	identify(conn, msg.Payload)
	// case messaging.MsgTypePing:
	// 	sendPong(conn)
	// default:
	// 	slog.Error("Unknown message type", slog.Int("type", int(msg.Type)))
	// }
	// }

	// return
	//
	// buf := make([]byte, 1024)
	// n, err := conn.Read(buf)
	// if err != nil {
	// 	slog.Error("Error reading Nickname", slog.Any("error", err))
	// 	return
	// }
	// clientName := string(buf[:n])
	//
	// clientsMutex.Lock()
	// clients[clientName] = Client{Nickname: clientName, Conn: conn}
	// slog.Info("Connected as", slog.String("Nickname", clientName))
	// logActiveClients()
	// clientsMutex.Unlock()
	//
	// for {
	// 	n, err := conn.Read(buf)
	// 	if err != nil {
	// 		if err == io.EOF {
	// 			clientsMutex.Lock()
	// 			delete(clients, clientName)
	// 			slog.Info("Client disconnected")
	// 			logActiveClients()
	// 			clientsMutex.Unlock()
	// 			break
	// 		}
	//
	// 		slog.Error("Error reading from connection", slog.Any("error", err))
	// 		return
	// 	}
	//
	// 	msg, err := messaging.Deserialize(buf[:n])
	// 	if err != nil {
	// 		slog.Error("Failed to deserialize ping", slog.Any("error", err))
	// 		return
	// 	}
	//
	// 	slog.Info("Received Ping", slog.Int64("timestamp", msg.Timestamp))
	//
	// 	err = sendPong(conn)
	// 	if err != nil {
	// 		slog.Error("Failed to send Pong", slog.Any("error", err))
	// 	}
	// }
}

func disconnectClient(conn net.Conn, nickname string) {
	clientsMutex.Lock()
	delete(clients, nickname)
	slog.Info("Client disconnected")
	logActiveClients()
	clientsMutex.Unlock()
}

// func sendPong(conn net.Conn) (int, error) {
// 	pongMessage := messaging.CreateMessage(messaging.MsgTypePong, nil)
// 	return conn.Write(pongMessage)
// }

// func handleIdentify(conn net.Conn, payload []byte) {
// 	nickname := string(payload)
// 	clientsMutex.Lock()
// 	defer clientsMutex.Unlock()
// 	clients[nickname] = Client{Nickname: nickname, Conn: conn}
// }

// func sendPong(conn net.Conn) error {
// 	timestamp := time.Now().Unix()
// 	pongMessage, err := messaging.Serialize(timestamp)
// 	if err != nil {
// 		return err
// 	}
//
// 	_, err = conn.Write(pongMessage)
// 	return err
// }

func logActiveClients() {
	slog.Info("Active clients:")
	for name := range clients {
		slog.Info("client", slog.String("Nickname", name))
	}
}

func generateConnectionID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
