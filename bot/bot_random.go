package bot

import (
	"math/rand"
	"time"

	"it1shka.com/checkers-server/gamelogic"
)

type BotRandom struct {
  rnd *rand.Rand
}

func InitBotRandom() BotRandom {
  currentTime := time.Now().Unix()
  source := rand.NewSource(currentTime)
  rnd := rand.New(source)
  return BotRandom{rnd}
}

func (bot BotRandom) Name() string {
  return "random"
}

func (bot BotRandom) Move(board gamelogic.Board) (gamelogic.BoardMove, bool) {
  availableMoves := board.CurrentPossibleMoves()
  if len(availableMoves) <= 0 {
    return gamelogic.BoardMove{}, false
  }
  choice := bot.rnd.Intn(len(availableMoves))
  move := availableMoves[choice]
  return move, true
}
