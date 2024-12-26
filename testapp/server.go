package testapp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"it1shka.com/checkers-server/bot"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func RunTestApp(port string) {
	fmt.Printf("Running test frontend on host: http://localhost%s\n", port)
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./testapp/webapp"))
	mux.Handle("/", fileServer)
	mux.HandleFunc("/ws-connect", handleWebsocket)
	mux.HandleFunc("/bot-names", handleBotNames)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal(err)
	}
}

func handleWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "failed to upgrade connection", http.StatusInternalServerError)
		return
	}
	defer conn.Close()
	botName := r.URL.Query().Get("bot")
	playerColor := r.URL.Query().Get("color")
	handleGame(conn, botName, playerColor)
}

func handleBotNames(w http.ResponseWriter, r *http.Request) {
	botNames := bot.GetBotNames()
	if err := json.NewEncoder(w).Encode(botNames); err != nil {
		http.Error(w, "serialization failure", http.StatusInternalServerError)
	}
}
