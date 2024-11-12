package testapp

import (
	"github.com/gorilla/websocket"
	"it1shka.com/checkers-server/bot"
	"it1shka.com/checkers-server/gamelogic"
)

type move struct {
	from gamelogic.PieceSquare
	to   gamelogic.PieceSquare
}

type player struct {
	conn        *websocket.Conn
	connFailure chan<- struct{}
	restart     chan<- struct{}
	turn        gamelogic.PieceColor
	board       <-chan gamelogic.Board
	moves       chan<- move
	winner      <-chan gamelogic.PieceColor
	done        <-chan struct{}
}

func (p *player) handle() {
	for {
		select {}
	}
}

type enemy struct {
	brain bot.Bot
	turn  gamelogic.PieceColor
	board <-chan gamelogic.Board
	moves chan<- move
	done  <-chan struct{}
}

func (e *enemy) handle() {
	for {
		select {}
	}
}

func handleGame(conn *websocket.Conn, botName, playerColor string) {
	playerTurn := gamelogic.BLACK
	if playerColor == "red" {
		playerTurn = gamelogic.RED
	}
	enemyBrain, ok := bot.GetBotByName(botName)
	if !ok {
		enemyBrain = bot.InitBotRandom()
	}

	playerBoard, enemyBoard := make(chan gamelogic.Board), make(chan gamelogic.Board)
	playerMoves, enemyMoves := make(chan move), make(chan move)

	connFailure := make(chan struct{})
	restart := make(chan struct{})
	winner := make(chan gamelogic.PieceColor)
	done := make(chan struct{})

	player := player{
		conn:        conn,
		connFailure: connFailure,
		restart:     restart,
		turn:        playerTurn,
		board:       playerBoard,
		moves:       playerMoves,
		winner:      winner,
		done:        done,
	}
	enemy := enemy{
		brain: enemyBrain,
		turn:  playerTurn.Opposite(),
		board: enemyBoard,
		moves: enemyMoves,
		done:  done,
	}

	go player.handle()
	go enemy.handle()

	var boardMemory map[string]bool
	var currentBoard gamelogic.Board

	startGame := func() {
		boardMemory = make(map[string]bool)
		currentBoard = gamelogic.InitBoard()
		boardMemory[currentBoard.String()] = true
		playerBoard <- currentBoard.Copy()
		enemyBoard <- currentBoard.Copy()
	}

	startGame()

	for {
		select {
		case <-connFailure:

		}
	}
}
