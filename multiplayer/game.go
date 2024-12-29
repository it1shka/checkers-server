package multiplayer

import (
	"it1shka.com/checkers-server/gamelogic"
)

type game struct {
	state       *gamelogic.GameSession
	playerBlack *player
	playerRed   *player
	moves       chan authoredMove
	done        chan bool
}

// TODO: red/black randomization
func newGame(playerBlack, playerRed *player) *game {
	return &game{
		state:       gamelogic.NewGameSession(),
		playerBlack: playerBlack,
		playerRed:   playerRed,
		moves:       make(chan authoredMove),
		done:        make(chan bool),
	}
}

func (g *game) startAsync() {
	go g.handlePlayer(g.playerBlack)
	go g.handlePlayer(g.playerRed)
}

func (g *game) handlePlayer(player *player) {
	for move := range player.movesChannel {
		g.moves <- authoredMove{
			author: player,
			move:   move,
		}
	}
}

func (g *game) retreat(player *player) {

}
