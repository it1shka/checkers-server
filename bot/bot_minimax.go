package bot

import (
	"fmt"

	"it1shka.com/checkers-server/gamelogic"
)

type BotMinimax struct {
  depth uint
}

func InitBotMinimax(depth uint) BotMinimax {
  return BotMinimax{depth}
}

func (bot BotMinimax) Name() string {
  return fmt.Sprintf("minimax-%d", bot.depth)
}

func (bot BotMinimax) Move(board gamelogic.Board) (gamelogic.BoardMove, bool) {
  // TODO: complete this function
  panic("Not implemented!")
}
