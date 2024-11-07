package testapp

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"it1shka.com/checkers-server/bot"
	"it1shka.com/checkers-server/gamelogic"
)

func RunLocally() {
	clearTerminal()
	bot := chooseBot()
	board := gamelogic.InitBoard()
	for {
		clearTerminal()
		fmt.Printf("Playing agains: %s\n", bot.Name())
		fmt.Println()
		fmt.Println(board)
		fmt.Println()
		availableMoves := board.CurrentPossibleMoves()
		if len(availableMoves) <= 0 {
			if board.Turn() == gamelogic.BLACK {
				fmt.Println("You lost")
			} else {
				fmt.Println("You won! Congrats!")
			}
			break
		}
		if board.Turn() == gamelogic.BLACK {
			from, ok := enterSquare("From")
			if !ok {
				continue
			}
			to, ok := enterSquare("To")
			if !ok {
				continue
			}
			nextBoard, ok := board.MakeMove(from, to)
			if ok {
				board = nextBoard
			}
			continue
		}
		move, ok := bot.Move(board)
		if !ok {
			fmt.Println("Bot failed. Aborting")
			break
		}
		nextBoard, ok := board.MakeMove(move.From, move.To)
		if ok {
			board = nextBoard
		}
	}
}

func enterSquare(label string) (gamelogic.PieceSquare, bool) {
	var position gamelogic.PiecePosition
	fmt.Printf("%s (row;col): ", label)
	if _, err := fmt.Scanf("%d;%d", &position.Row, &position.Column); err != nil {
		return -1, false
	}
	if !position.IsValid() {
		return -1, false
	}
	return position.ToSquare(), true
}

func chooseBot() bot.Bot {
	botNames := bot.GetBotNames()
	for {
		fmt.Println("Available Bots:")
		for index, botName := range botNames {
			fmt.Printf("%d) %s\n", index+1, botName)
		}
		fmt.Printf("Choose bot by entering a number (%d-%d):\n", 1, len(botNames))
		var choice int
		if _, err := fmt.Scanf("%d", &choice); err != nil {
			fmt.Println("Not a number. Please, try again!")
			fmt.Scanln()
			continue
		}
		if choice < 1 || choice > len(botNames) {
			fmt.Println("Not in range. Please, try again!")
			continue
		}
		bot, ok := bot.GetBotByName(botNames[choice-1])
		if !ok {
			fmt.Println("Please, try again!")
			continue
		}
		return bot
	}
}

func clearTerminal() bool {
	system := runtime.GOOS
	var cmd *exec.Cmd = nil
	if system == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	if cmd == nil {
		return false
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
	return cmd.Err == nil
}
