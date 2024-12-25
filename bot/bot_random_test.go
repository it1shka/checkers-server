package bot_test

import (
	"testing"

	"it1shka.com/checkers-server/bot"
	"it1shka.com/checkers-server/gamelogic"
)

func TestBotRandomCorrectness(t *testing.T) {
	botRandom := bot.InitBotRandom()
	board := gamelogic.InitBoard()

	for len(board.CurrentPossibleMoves()) > 0 {
		move, ok := botRandom.Move(board)
		if !ok {
			t.Fatal("bot failed to move")
		}
		nextBoard, ok := board.MakeMove(move.From, move.To)
		if !ok {
			t.Fatal("board failed to update")
		}
		board = nextBoard
	}
}
