package multiplayer

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	SKT_MSG_UUID  = "uuid"
	SKT_MSG_CRASH = "crash"
)

type SocketMessage struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

func (msg SocketMessage) SendTo(conn *websocket.Conn) error {
	return conn.WriteJSON(msg)
}

const (
	CLT_MSG_JOIN  = "join"
	CLT_MSG_LEAVE = "leave"
)

type ClientMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

func ExtractPayload[T any](conn *websocket.Conn, msg ClientMessage) (T, bool) {
	var output T
	if err := json.Unmarshal(msg.Payload, &output); err == nil {
		return output, true
	}
	MessageCrash("received malformed message").SendTo(conn)
	return output, false
}

// catalog of messages to send

func MessageUUID(id uuid.UUID) SocketMessage {
	return SocketMessage{
		Type:    SKT_MSG_UUID,
		Payload: id.String(),
	}
}

func MessageCrash(reason string) SocketMessage {
	return SocketMessage{
		Type:    SKT_MSG_CRASH,
		Payload: reason,
	}
}
