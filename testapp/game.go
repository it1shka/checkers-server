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
	board       <-chan gamelogic.Board
	moves       chan<- move
	winner      <-chan string
	done        <-chan struct{}
}

func (p *player) handleMessages() {
PlayerMessagesForSelect:
	for {
		select {
		case <-p.done:
			break PlayerMessagesForSelect
		default:
			var msg map[string]any
			err := p.conn.ReadJSON(&msg)
			if err != nil {
				p.connFailure <- struct{}{}
			} else {
				p.handleMessage(msg)
			}
		}
	}
}

func (p *player) handleMessage(msg map[string]any) {
	rawTag, ok := msg["tag"]
	if !ok {
		return
	}
	tag, ok := rawTag.(string)
	if !ok {
		return
	}
	switch tag {
	case "restart":
		p.restart <- struct{}{}
	case "move":
		rawFrom, ok := msg["from"]
		if !ok {
			return
		}
		from, ok := rawFrom.(int)
		if !ok {
			return
		}
		rawTo, ok := msg["to"]
		if !ok {
			return
		}
		to, ok := rawTo.(int)
		if !ok {
			return
		}
		fromSquare := gamelogic.PieceSquare(from)
		toSquare := gamelogic.PieceSquare(to)
		if !fromSquare.IsValid() || !toSquare.IsValid() {
			return
		}
		p.moves <- move{fromSquare, toSquare}
	}
}

func (p *player) handle() {
	go p.handleMessages()

PlayerForSelect:
	for {
		select {
		case brd := <-p.board:
			turn := "black"
			if brd.Turn() == gamelogic.RED {
				turn = "red"
			}
			var pieces []map[string]any
			for _, piece := range brd.Pieces() {
				pieceColor := "black"
				if piece.Color == gamelogic.RED {
					pieceColor = "red"
				}
				pieceType := "man"
				if piece.Type == gamelogic.KING {
					pieceType = "king"
				}
				pieceSquare := int(piece.Square)
				pieces = append(pieces, map[string]any{
					"color":  pieceColor,
					"type":   pieceType,
					"square": pieceSquare,
				})
			}
			msg := map[string]any{
				"tag":    "board",
				"turn":   turn,
				"pieces": pieces,
			}
			p.conn.WriteJSON(msg)
		case wnr := <-p.winner:
			msg := map[string]any{
				"tag":    "winner",
				"winner": wnr,
			}
			p.conn.WriteJSON(msg)
		case <-p.done:
			break PlayerForSelect
		}
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
EnemyForSelect:
	for {
		select {
		case brd := <-e.board:
			if brd.Turn() == e.turn {
				if mv, ok := e.brain.Move(brd); ok {
					e.moves <- move{mv.From, mv.To}
				}
			}
		case <-e.done:
			break EnemyForSelect
		}
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
	winner := make(chan string)
	done := make(chan struct{})

	player := player{
		conn:        conn,
		connFailure: connFailure,
		restart:     restart,
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

	var boardMemory map[string]int
	var currentBoard gamelogic.Board

	startGame := func() {
		boardMemory = make(map[string]int)
		currentBoard = gamelogic.InitBoard()
		boardMemory[currentBoard.String()] = 1
		playerBoard <- currentBoard.Copy()
		enemyBoard <- currentBoard.Copy()
	}

	makeMove := func(mv move) {
		nextBoard, ok := currentBoard.MakeMove(mv.from, mv.to)
		if !ok {
			return
		}
		currentBoard = nextBoard
		boardMemory[nextBoard.String()]++
		playerBoard <- nextBoard.Copy()
		enemyBoard <- nextBoard.Copy()

		if len(nextBoard.CurrentPossibleMoves()) <= 0 {
			wnr := nextBoard.Turn().Opposite()
			if wnr == gamelogic.BLACK {
				winner <- "black"
			} else {
				winner <- "red"
			}
			return
		}

		if boardMemory[nextBoard.String()] >= 3 {
			winner <- "draw"
		}
	}

	startGame()

PrimaryForSelect:
	for {
		select {
		case <-connFailure:
			close(done)
			break PrimaryForSelect
		case <-restart:
			startGame()
		case <-done:
			break PrimaryForSelect
		case mv := <-playerMoves:
			if currentBoard.Turn() == playerTurn {
				makeMove(mv)
			}
		case mv := <-enemyMoves:
			if currentBoard.Turn() == playerTurn.Opposite() {
				makeMove(mv)
			}
		}
	}
}
