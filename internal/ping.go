package ping

import (
	"bytes"
	"encoding/binary"
)

type PingMessage struct {
	Type      uint8
	Timestamp int64
}

func Serialize(timestamp int64) ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.LittleEndian, uint8(0x01)); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.LittleEndian, timestamp); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func Deserialize(data []byte) (PingMessage, error) {
	buf := bytes.NewReader(data)
	var msg PingMessage

	if err := binary.Read(buf, binary.LittleEndian, &msg.Type); err != nil {
		return msg, err
	}

	if err := binary.Read(buf, binary.LittleEndian, &msg.Timestamp); err != nil {
		return msg, err
	}

	return msg, nil
}
