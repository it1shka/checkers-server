package multiplayer

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type SocketMessage struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

func (msg SocketMessage) SendTo(conn *websocket.Conn) error {
	return conn.WriteJSON(msg)
}

// catalog of messages to send

func MessageUUID(id uuid.UUID) SocketMessage {
	return SocketMessage{
		Type:    "uuid",
		Payload: id.String(),
	}
}

func MessageCrash(reason string) SocketMessage {
	return SocketMessage{
		Type:    "crash",
		Payload: reason,
	}
}
