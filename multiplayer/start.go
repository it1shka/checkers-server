package multiplayer

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var server = NewServer()

func StartMultiplayer(port string) {
	fmt.Printf("Starting up multiplayer server on: %s\n", port)
	http.HandleFunc("/ws-connect", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			http.Error(w, "failed to upgrade connection", http.StatusInternalServerError)
			return
		}
		go server.HandleConnection(conn)
	})
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
