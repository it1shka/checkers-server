package multiplayer

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type Server struct {
	uuids      UUIDs
	dispatcher map[string]func(*websocket.Conn, json.RawMessage)
}

func NewServer() *Server {
	server := &Server{
		uuids:      InitUUIDs(),
		dispatcher: nil,
	}
	server.dispatcher = map[string]func(*websocket.Conn, json.RawMessage){
		CLT_MSG_JOIN:  server.handleClientJoin,
		CLT_MSG_LEAVE: server.handleClientLeave,
	}
	return server
}

func (server *Server) HandleConnection(conn *websocket.Conn) {
	defer conn.Close()
	id, ok := server.uuids.GenerateFor(conn)
	if !ok {
		MessageCrash("failed to register connection").SendTo(conn)
		return
	}
	defer server.uuids.DeactivateFor(conn)
	MessageUUID(id).SendTo(conn)

	for {
		var message ClientMessage
		if err := conn.ReadJSON(&message); err != nil {
			MessageCrash("failed to read message").SendTo(conn)
			continue
		}
		if handler, ok := server.dispatcher[message.Type]; ok {
			handler(conn, message.Payload)
		} else {
			MessageCrash("received unknown message type").SendTo(conn)
		}
	}
}

func (server *Server) handleClientJoin(conn *websocket.Conn, payload json.RawMessage) {
	// TODO: request to enter the queue
}

func (server *Server) handleClientLeave(conn *websocket.Conn, payload json.RawMessage) {
	// TODO: request to leave the queue
}
