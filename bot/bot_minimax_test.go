package bot_test

import (
	"testing"

	"it1shka.com/checkers-server/bot"
	"it1shka.com/checkers-server/gamelogic"
)

const ACCURACY_ITERATIONS = 20

func TestBotAccuracy(t *testing.T) {
	botBlack := bot.InitBotMinimax(6)
	botRandom := bot.InitBotRandom()

	wins := 0
	for i := 0; i < ACCURACY_ITERATIONS; i++ {
		board := gamelogic.InitBoard()
		for len(board.CurrentPossibleMoves()) > 0 {
			var move gamelogic.BoardMove
			var ok bool
			if board.Turn() == gamelogic.BLACK {
				move, ok = botBlack.Move(board)
			} else {
				move, ok = botRandom.Move(board)
			}
			if !ok {
				t.Fatal("bot failed to move")
			}
			nextBoard, ok := board.MakeMove(move.From, move.To)
			if !ok {
				t.Fatal("board failed to update")
			}
			board = nextBoard
		}
		if board.Turn() == gamelogic.RED {
			wins++
			t.Logf("won %d time(s)", wins)
		}
	}

	loses := ACCURACY_ITERATIONS - wins
	if loses > 0 {
		t.Fatalf("bot failed %d time(s)", loses)
	}
}

const LEVELING_ITERATIONS = 5
const MOVES_LIMIT = 100

func TestBotDifficultyLeveling(t *testing.T) {
	botBlack := bot.InitBotMinimax(8)
	botRed := bot.InitBotMinimax(4)

	wins := 0
	skipped := 0
	for i := 0; i < LEVELING_ITERATIONS; i++ {
		board := gamelogic.InitBoard()
		localMoves := 0
		for len(board.CurrentPossibleMoves()) > 0 && localMoves < MOVES_LIMIT {
			var move gamelogic.BoardMove
			var ok bool
			if board.Turn() == gamelogic.BLACK {
				move, ok = botBlack.Move(board)
			} else {
				move, ok = botRed.Move(board)
			}
			if !ok {
				t.Fatal("bot failed to move")
			}
			nextBoard, ok := board.MakeMove(move.From, move.To)
			if !ok {
				t.Fatal("board failed to update")
			}
			board = nextBoard
			localMoves++
		}
		if len(board.CurrentPossibleMoves()) > 0 {
			skipped++
			t.Logf("skipped")
			continue
		}
		if board.Turn() == gamelogic.RED {
			wins++
			t.Logf("advanced bot won %d time(s)", wins)
		} else {
			t.Log("advanced bot failed")
		}
	}

	loses := LEVELING_ITERATIONS - wins - skipped
	if loses > 0 {
		t.Fatalf("advanced bot failed %d time(s)", loses)
	}
}
