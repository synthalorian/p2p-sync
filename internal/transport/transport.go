package transport

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
)

type MessageType string

const (
	MsgHello       MessageType = "hello"
	MsgFileList    MessageType = "file_list"
	MsgFileRequest MessageType = "file_request"
	MsgFileData    MessageType = "file_data"
	MsgFileDone    MessageType = "file_done"
)

type Message struct {
	Type    MessageType     `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type HelloPayload struct {
	Name string `json:"name"`
}

type FileInfo struct {
	Name string `json:"name"`
	Hash string `json:"hash"`
	Size int64  `json:"size"`
}

type FileListPayload struct {
	Files []FileInfo `json:"files"`
}

type FileRequestPayload struct {
	Name string `json:"name"`
}

type FileDataPayload struct {
	Name string `json:"name"`
	Data []byte `json:"data"`
	Size int64  `json:"size"`
}

type FileDonePayload struct {
	Name string `json:"name"`
}

func ReadMessage(conn net.Conn) (*Message, error) {
	var length uint32
	if err := binary.Read(conn, binary.BigEndian, &length); err != nil {
		return nil, fmt.Errorf("read length: %w", err)
	}
	buf := make([]byte, length)
	if _, err := io.ReadFull(conn, buf); err != nil {
		return nil, fmt.Errorf("read payload: %w", err)
	}
	var msg Message
	if err := json.Unmarshal(buf, &msg); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	return &msg, nil
}

func WriteMessage(conn net.Conn, msg *Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	if err := binary.Write(conn, binary.BigEndian, uint32(len(data))); err != nil {
		return err
	}
	_, err = conn.Write(data)
	return err
}

func SendPayload(conn net.Conn, msgType MessageType, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	msg := &Message{
		Type:    msgType,
		Payload: json.RawMessage(data),
	}
	return WriteMessage(conn, msg)
}
