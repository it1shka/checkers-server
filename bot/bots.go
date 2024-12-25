package bot

import "it1shka.com/checkers-server/gamelogic"

type Bot interface {
	Name() string
	Move(board gamelogic.Board) (gamelogic.BoardMove, bool)
}

func GetBots() []Bot {
	return []Bot{
		InitBotRandom(),
		InitBotMinimax(2),
		InitBotMinimax(4),
		InitBotMinimax(6),
		InitBotMinimax(10),
	}
}

func GetBotNames() []string {
	bots := GetBots()
	var output []string
	for _, bot := range bots {
		name := bot.Name()
		output = append(output, name)
	}
	return output
}

func GetBotByName(name string) (Bot, bool) {
	bots := GetBots()
	for _, bot := range bots {
		if bot.Name() == name {
			return bot, true
		}
	}
	return nil, false
}
