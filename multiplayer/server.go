package multiplayer

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/gorilla/websocket"
)

const serverReadBufferSize = 1024
const serverWriteBufferSize = 1024

type Server struct {
	upgrader websocket.Upgrader
	decoder  *schema.Decoder
}

func NewServer() *Server {
	return &Server{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  serverReadBufferSize,
			WriteBufferSize: serverWriteBufferSize,
		},
		decoder: schema.NewDecoder(),
	}
}

func (s *Server) Start(port string) {
	fmt.Printf("Running multiplayer server on port: %s\n", port)
	mux := http.NewServeMux()
	mux.HandleFunc("/ws-connect", s.handleRequest)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalln(err)
	}
}

func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	var playerInfo playerInfo
	if err := s.decoder.Decode(&playerInfo, r.URL.Query()); err != nil {
		msg := "please, provide nickname, rating and region"
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		msg := "disconnected due to unknown reason"
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	player, err := newPlayer(conn, playerInfo)
	if err != nil {
		log.Println(err)
		msg := "disconnected due to unknown reason"
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	player.startAsync()
}
