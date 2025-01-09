package multiplayer

import (
	"log"
	"math/rand"
	"sync/atomic"
	"time"

	"it1shka.com/checkers-server/bot"
	"it1shka.com/checkers-server/gamelogic"
)

type pseudogame struct {
	state      *gamelogic.GameSession
	human      *player
	humanColor gamelogic.PieceColor
	bot        bot.Bot
	timeBlack  atomic.Int32
	timeRed    atomic.Int32
	timeChange chan bool
	timeFlag   atomic.Bool
	moves      chan authoredMove
	done       chan bool
}

func newPseudogame(human *player) *pseudogame {
	guess := randomBool()
	var humanColor gamelogic.PieceColor
	if guess {
		humanColor = gamelogic.BLACK
	} else {
		humanColor = gamelogic.RED
	}

	bots := bot.GetBots()
	chosenBot := bots[rand.Intn(len(bots))]

	output := &pseudogame{
		state:      gamelogic.NewGameSession(),
		human:      human,
		humanColor: humanColor,
		bot:        chosenBot,
		timeFlag:   atomic.Bool{},
		timeChange: make(chan bool),
		timeBlack:  atomic.Int32{},
		timeRed:    atomic.Int32{},
		moves:      make(chan authoredMove),
		done:       make(chan bool),
	}
	output.timeBlack.Store(maxGameTime)
	output.timeRed.Store(maxGameTime)
	output.timeFlag.Store(timeBlackFlag)
	return output
}

func (g *pseudogame) startAsync() {
	if g.humanColor == gamelogic.BLACK {
		g.human.sendMessage(getOutMsgColor(colorBlack))
	} else {
		g.human.sendMessage(getOutMsgColor(colorRed))
	}
	pseudoEnemyInfo := getPseudoPlayerInfo()
	g.human.sendMessage(getOutMsgPseudoEnemy(pseudoEnemyInfo))

	boardMsg := getOutMsgBoard(g.state.Board())
	statusMsg := getOutMsgStatus(g.state.Status())
	selfTimeMsg := getOutMsgTime(g.humanColor.String(), maxGameTime)
	enemyTimeMsg := getOutMsgTime(g.humanColor.Opposite().String(), maxGameTime)
	g.human.sendMessage(boardMsg)
	g.human.sendMessage(statusMsg)
	g.human.sendMessage(selfTimeMsg)
	g.human.sendMessage(enemyTimeMsg)

	go func() {
	ListeningMoves:
		for {
			select {
			case <-g.done:
				break ListeningMoves
			case move := <-g.moves:
				var moveColor gamelogic.PieceColor
				if move.author == g.human {
					moveColor = g.humanColor
				} else {
					moveColor = g.humanColor.Opposite()
				}
				previousTurn := g.state.Board().Turn()
				success := g.state.MakeMove(moveColor, move.move.from, move.move.to)
				if !success {
					continue ListeningMoves
				}
				boardMsg := getOutMsgBoard(g.state.Board())
				statusMsg := getOutMsgStatus(g.state.Status())
				g.human.sendMessage(boardMsg)
				g.human.sendMessage(statusMsg)

				if g.state.Status() != gamelogic.ACTIVE {
					g.finish()
					break ListeningMoves
				}
				if g.state.Board().Turn() != previousTurn {
					g.timeChange <- true
					g.timeFlag.Store(!g.timeFlag.Load())
					g.startClockAsync()
				}

				if g.state.Board().Turn() == g.humanColor.Opposite() {
					g.triggerBotAsync()
				}
			}
		}
	}()

	g.startClockAsync()

	if g.humanColor == gamelogic.RED {
		g.triggerBotAsync()
	}
}

func (g *pseudogame) triggerBotAsync() {
	go func() {
		move, ok := g.bot.Move(g.state.Board())
		if !ok {
			log.Println("bot failed")
			g.finish()
			return
		}
		authored := authoredMove{
			author: nil,
			move: incomingMove{
				from: move.From,
				to:   move.To,
			},
		}
		g.pushMove(authored)
	}()
}

func (g *pseudogame) startClockAsync() {
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
					g.human.sendMessage(timeMsg)
					if newTime <= 0 {
						statusMsg := getOutMsgStatus(gamelogic.RED_WIN)
						g.human.sendMessage(statusMsg)
						g.finish()
						break Ticking
					}
				} else {
					newTime := g.timeRed.Add(-1)
					timeMsg := getOutMsgTime(colorRed, newTime)
					g.human.sendMessage(timeMsg)
					if newTime <= 0 {
						statusMsg := getOutMsgStatus(gamelogic.BLACK_WIN)
						g.human.sendMessage(statusMsg)
						g.finish()
						break Ticking
					}
				}
			}
		}
	}()
}

func (g *pseudogame) retreat(loser *player) {
	// TODO: maybe send the message to human that he lost
	// TODO: in the future
	g.finish()
}

func (g *pseudogame) pushMove(move authoredMove) {
	select {
	case <-g.done:
		return
	default:
		g.moves <- move
	}
}

func (g *pseudogame) finish() {
	close(g.done)
	close(g.moves)
	close(g.timeChange)
	g.state = nil
	g.human = nil
}
