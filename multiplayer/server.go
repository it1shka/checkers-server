package multiplayer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/gorilla/websocket"
)

const SERVER_READ_BUFFER_SIZE = 1024
const SERVER_WRITE_BUFFER_SIZE = 1024

type Server struct {
	upgrader websocket.Upgrader
	decoder  *schema.Decoder
}

func NewServer() *Server {
	server := &Server{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  SERVER_READ_BUFFER_SIZE,
			WriteBufferSize: SERVER_WRITE_BUFFER_SIZE,
		},
		decoder: schema.NewDecoder(),
	}
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
	var metadata struct {
		Nickname string `schema:"nickname"`
		Rating   uint   `schema:"rating"`
		Region   string `schema:"region"`
	}

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

	defer conn.Close()

	for {
		var message struct {
			Type    string          `json:"type"`
			Payload json.RawMessage `json:"payload"`
		}
		if err := conn.ReadJSON(&message); err != nil {
			if !websocket.IsCloseError(err) && !websocket.IsUnexpectedCloseError(err) {
				log.Println(err)
			}
			return
		}
		switch message.Type {

		}
	}
}
