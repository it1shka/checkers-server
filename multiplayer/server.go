package multiplayer

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/schema"
	"github.com/gorilla/websocket"
)

// Helper records

type metadata struct {
	ID_       string
	Nickname_ string `schema:"nickname"`
	Rating_   uint   `schema:"rating"`
	Region_   string `schema:"region"`
}

func (m metadata) ID() string {
	return m.ID_
}
func (m metadata) Nickname() string {
	return m.Nickname_
}
func (m metadata) Rating() uint {
	return m.Rating_
}
func (m metadata) Region() string {
	return m.Region_
}

type message struct {
	Type_    string `json:"type"`
	Payload_ any    `json:"payload"`
}

func (m message) Type() string {
	return m.Type_
}
func (m message) Payload() any {
	return m.Payload_
}

// Server implementation

const SERVER_READ_BUFFER_SIZE = 1024
const SERVER_WRITE_BUFFER_SIZE = 1024
const MATCHMAKING_PERIOD = 5 * time.Second

type Server struct {
	upgrader    websocket.Upgrader
	decoder     *schema.Decoder
	multiplayer *Multiplayer
}

func NewServer() *Server {
	mult := NewMultiplayer()
	server := &Server{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  SERVER_READ_BUFFER_SIZE,
			WriteBufferSize: SERVER_WRITE_BUFFER_SIZE,
		},
		decoder:     schema.NewDecoder(),
		multiplayer: mult,
	}
	mult.StartMatchmaking(MATCHMAKING_PERIOD)
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

func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	var metadata metadata
	if err := s.decoder.Decode(&metadata, r.URL.Query()); err != nil {
		log.Println(err)
		http.Error(w, "failed to parse your request", http.StatusBadRequest)
		return
	}
	metadata.ID_ = uuid.NewString()

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "failed to upgrade connection", http.StatusMethodNotAllowed)
		return
	}
	defer conn.Close()

	s.multiplayer.RegisterPlayer(metadata, func(update IUpdate) {
		conn.WriteJSON(message{
			Type_:    update.Type(),
			Payload_: update.Payload(),
		})
	})
	defer s.multiplayer.UnregisterPlayer(metadata.ID_)

	for {
		var message message
		if err := conn.ReadJSON(&message); err != nil {
			if !websocket.IsCloseError(err) && !websocket.IsUnexpectedCloseError(err) {
				log.Println(err)
			}
			return
		}
		s.multiplayer.PushUpdate(metadata.ID_, message)
	}
}
