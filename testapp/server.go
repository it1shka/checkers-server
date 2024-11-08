package testapp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"it1shka.com/checkers-server/gamelogic"
)

const (
	BLACK = "black"
	RED   = "red"
	MAN   = "man"
	KING  = "king"
)

func RunTestApp() {
	fmt.Println("Running on: http://localhost:3333/")
	fileServer := http.FileServer(http.Dir("./testapp/webapp"))
	http.Handle("/", fileServer)
	http.HandleFunc("/api/start-game", handleStartGame)
	err := http.ListenAndServe(":3333", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func writeBoard(w http.ResponseWriter, board gamelogic.Board) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]any{
		"pieces": board.Pieces(),
		"status": "active",
	})
	if err != nil {
		message := "Serialization error"
		http.Error(w, message, http.StatusInternalServerError)
	}
}

func handleStartGame(w http.ResponseWriter, r *http.Request) {
	board := gamelogic.InitBoard()
	writeBoard(w, board)
}
