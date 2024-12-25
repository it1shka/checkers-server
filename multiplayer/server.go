package multiplayer

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/gorilla/websocket"
)

const SERVER_READ_BUFFER_SIZE = 1024
const SERVER_WRITE_BUFFER_SIZE = 1024

type Server struct {
	port     string
	upgrader websocket.Upgrader
	decoder  *schema.Decoder
	players  *PlayerCollection
}

// Note: to start server, use the following:
// NewServer("<your port>").Start()
func NewServer(port string) *Server {
	return &Server{
		port: port,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  SERVER_READ_BUFFER_SIZE,
			WriteBufferSize: SERVER_WRITE_BUFFER_SIZE,
		},
		decoder: schema.NewDecoder(),
		players: NewPlayerCollection(),
	}
}

func (s *Server) GetPort() string {
	return s.port
}

func (s *Server) Start() {
	fmt.Printf("Running multiplayer on: %s\n", s.port)
	mux := http.NewServeMux()
	mux.HandleFunc("/ws-connect", s.handleRequest)
	if err := http.ListenAndServe(s.port, mux); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	var player Player
	if err := s.decoder.Decode(&player, r.URL.Query()); err != nil {
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

	player.Conn = conn

	// TODO: nickname and rating should be
	// TODO: pulled out of the database -- never trust the frontend

	// TODO: handle connection
	fmt.Printf("%v", player)
	conn.Close()
}
