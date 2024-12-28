package multiplayer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/gorilla/websocket"
	"it1shka.com/checkers-server/gamelogic"
)

// Additional record types

type RequestMetadata struct {
	Nickname string `schema:"nickname"`
	Rating   uint   `schema:"rating"`
	Region   string `schema:"region"`
}

type IncomingMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

const (
	MSG_IN_JOIN  = "join"
	MSG_IN_LEAVE = "leave"
	MSG_IN_MOVE  = "move"
)

type OutcomingMessage struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

// Server implementation

const SERVER_READ_BUFFER_SIZE = 1024
const SERVER_WRITE_BUFFER_SIZE = 1024

type Server struct {
	upgrader    websocket.Upgrader
	decoder     *schema.Decoder
	matchmaking *Matchmaking[Player]
}

func NewServer() *Server {
	server := &Server{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  SERVER_READ_BUFFER_SIZE,
			WriteBufferSize: SERVER_WRITE_BUFFER_SIZE,
		},
		decoder:     schema.NewDecoder(),
		matchmaking: nil,
	}
	server.matchmaking = NewMatchmaking(server.handleMatch)
	return server
}

func (s *Server) Start(port string) {
	fmt.Printf("Running multiplayer on port: %s\n", port)
	mux := http.NewServeMux()
	mux.HandleFunc("/ws-connect", s.handleRequest)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) handleMatch(playerA, playerB Player) {
	go func() {
		playerA.NotifyAboutEnemy(playerB.Info())
		playerB.NotifyAboutEnemy(playerA.Info())
		session := NewSession(playerA.ID(), playerB.ID())
		updateHandler := func(update Update) {
			playerA.NotifyAboutUpdate(update)
			playerB.NotifyAboutUpdate(update)
			if update.Status != gamelogic.ACTIVE {
				session.Stop()
				playerA.Stop()
				playerB.Stop()
			}
		}

	}()
}

func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	var metadata RequestMetadata
	if err := s.decoder.Decode(&metadata, r.URL.Query()); err != nil {
		log.Println(err)
		http.Error(w, "failed to parse your request", http.StatusBadRequest)
		return
	}
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "failed to upgrade connection", http.StatusMethodNotAllowed)
		return
	}

	go s.handleConnection(conn, metadata)
}

func (s *Server) handleConnection(conn *websocket.Conn, meta RequestMetadata) {
	defer conn.Close()
	for {
		var message IncomingMessage
		if err := conn.ReadJSON(&message); err != nil {
			if !websocket.IsCloseError(err) && !websocket.IsUnexpectedCloseError(err) {
				log.Println(err)
			}
			return
		}

		switch message.Type {
		case MSG_IN_JOIN:

		case MSG_IN_LEAVE:

		case MSG_IN_MOVE:

		default:

		}
	}
}
