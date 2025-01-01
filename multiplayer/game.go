package multiplayer

import (
	"math/rand"
	"sync/atomic"
	"time"

	"it1shka.com/checkers-server/gamelogic"
)

const maxGameTime int32 = 300

const (
	timeBlackFlag = true
	timeRedFlag   = false
)

const (
	colorBlack = "black"
	colorRed   = "red"
)

type game struct {
	state       *gamelogic.GameSession
	playerBlack *player
	playerRed   *player
	timeBlack   atomic.Int32
	timeRed     atomic.Int32
	timeChange  chan bool
	timeFlag    atomic.Bool
	moves       chan authoredMove
	done        chan bool
}

func newGame(playerA, playerB *player) *game {
	var playerBlack, playerRed *player
	if rand.Intn(2) == 0 {
		playerBlack, playerRed = playerA, playerB
	} else {
		playerBlack, playerRed = playerB, playerA
	}
	output := &game{
		state:       gamelogic.NewGameSession(),
		playerBlack: playerBlack,
		playerRed:   playerRed,
		timeFlag:    atomic.Bool{},
		timeChange:  make(chan bool),
		timeBlack:   atomic.Int32{},
		timeRed:     atomic.Int32{},
		moves:       make(chan authoredMove),
		done:        make(chan bool),
	}
	output.timeBlack.Store(maxGameTime)
	output.timeRed.Store(maxGameTime)
	output.timeFlag.Store(timeBlackFlag)
	return output
}

func (g *game) colorOf(player *player) gamelogic.PieceColor {
	if g.playerBlack == player {
		return gamelogic.BLACK
	}
	return gamelogic.RED
}

func (g *game) oppositeOf(player *player) *player {
	if g.playerBlack == player {
		return g.playerRed
	}
	return g.playerBlack
}

func (g *game) startAsync() {
	g.playerBlack.sendMessage(getOutMsgColor(colorBlack))
	g.playerRed.sendMessage(getOutMsgColor(colorRed))

	g.playerBlack.sendMessage(getOutMsgEnemy(g.playerRed))
	g.playerRed.sendMessage(getOutMsgEnemy(g.playerBlack))

	go func() {
	ListeningMoves:
		for {
			select {
			case <-g.done:
				break ListeningMoves
			case move := <-g.moves:
				color := g.colorOf(move.author)
				success := g.state.MakeMove(color, move.move.from, move.move.to)
				if !success {
					continue ListeningMoves
				}
				boardMsg := getOutMsgBoard(g.state.Board())
				statusMsg := getOutMsgStatus(g.state.Status())
				for _, player := range []*player{g.playerBlack, g.playerRed} {
					player.sendMessage(boardMsg)
					player.sendMessage(statusMsg)
				}

				if g.state.Status() != gamelogic.ACTIVE {
					g.finish()
					break ListeningMoves
				}

				g.timeChange <- true
				g.startClockAsync()
			}
		}
	}()

	g.startClockAsync()
}

func (g *game) startClockAsync() {
	go func() {
		ticker := time.Tick(time.Second)
		flag := g.timeFlag.Load()
	Ticking:
		for {
			select {
			case <-g.done:
				break Ticking
			case <-g.timeChange:
				break Ticking
			case <-ticker:
				if flag == timeBlackFlag {
					newTime := g.timeBlack.Add(-1)
					timeMsg := getOutMsgTime(colorBlack, newTime)
					g.playerBlack.sendMessage(timeMsg)
					g.playerRed.sendMessage(timeMsg)
					if newTime <= 0 {
						g.retreat(g.playerBlack)
						break Ticking
					}
				} else {
					newTime := g.timeRed.Add(-1)
					timeMsg := getOutMsgTime(colorRed, newTime)
					g.playerBlack.sendMessage(timeMsg)
					g.playerRed.sendMessage(timeMsg)
					if newTime <= 0 {
						g.retreat(g.playerRed)
						break Ticking
					}
				}
			}
		}
	}()
}

func (g *game) retreat(loser *player) {
	defer g.finish()
	winner := g.oppositeOf(loser)
	var status gamelogic.GameStatus
	if winner == g.playerBlack {
		status = gamelogic.BLACK_WIN
	} else {
		status = gamelogic.RED_WIN
	}
	statusMsg := getOutMsgStatus(status)
	for _, player := range []*player{g.playerBlack, g.playerRed} {
		player.sendMessage(statusMsg)
		player.sendMessage(statusMsg)
	}
}

func (g *game) pushMove(move authoredMove) {
	select {
	case <-g.done:
		return
	default:
		g.moves <- move
	}
}

func (g *game) finish() {
	close(g.done)
	close(g.moves)
	close(g.timeChange)
	g.state = nil
	g.playerBlack = nil
	g.playerRed = nil
}
