package messaging

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

func SendMessage(conn net.Conn, msg Message) error {
	data, err := msg.Serialize()
	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	return err
}

func ReceiveMessage(conn net.Conn) (Message, error) {
	typeBuf := make([]byte, 1)
	_, err := conn.Read(typeBuf)
	if err != nil {
		return nil, err
	}

	messageType := typeBuf[0]

	var msg Message
	switch messageType {
	case 1:
		msg = &ConnStartMessage{}
	case 2:
		msg = &PingMessage{}
	default:
		return nil, fmt.Errorf("unknown message type: %d", messageType)
	}

	dataBuf := make([]byte, 1024)
	n, err := conn.Read(dataBuf)
	if err != nil {
		return nil, err
	}

	err = msg.Deserialize(dataBuf[:n])
	return msg, err
}

type Message interface {
	MessageType() uint8
	Serialize() ([]byte, error)
	Deserialize([]byte) error
}

type ConnStartMessage struct {
	Id string
}

func (m *ConnStartMessage) MessageType() uint8 {
	return 1
}

func (m *ConnStartMessage) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := buf.WriteByte(m.MessageType()); err != nil {
		return nil, err
	}

	if err := writeString(buf, m.Id); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (m *ConnStartMessage) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)

	id, err := readString(buf)
	if err != nil {
		return err
	}
	m.Id = id

	return nil
}

type PingMessage struct {
	Timestamp int64
}

func (m *PingMessage) MessageType() uint8 {
	return 2
}

func (m *PingMessage) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := buf.WriteByte(m.MessageType()); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.LittleEndian, m.Timestamp); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (m *PingMessage) Deserialize(data []byte) error {
	buf := bytes.NewReader(data)

	err := binary.Read(buf, binary.LittleEndian, &m.Timestamp)

	return err
}

func writeString(buf *bytes.Buffer, s string) error {
	length := uint16(len(s))
	if err := binary.Write(buf, binary.LittleEndian, length); err != nil {
		return err
	}
	if _, err := buf.Write([]byte(s)); err != nil {
		return err
	}
	return nil
}

func readString(buf *bytes.Reader) (string, error) {
	var length uint16
	if err := binary.Read(buf, binary.LittleEndian, &length); err != nil {
		return "", err
	}
	strBytes := make([]byte, length)
	if _, err := buf.Read(strBytes); err != nil {
		return "", err
	}
	return string(strBytes), nil
}
