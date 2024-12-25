package multiplayer

import "github.com/gorilla/websocket"

type Server struct {
	uuids UUIDs
}

func NewServer() *Server {
	return &Server{
		uuids: InitUUIDs(),
	}
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

	// TODO: ...
}
