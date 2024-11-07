package bot

import (
	"fmt"

	"it1shka.com/checkers-server/gamelogic"
)

func getCenterSquares() []gamelogic.PieceSquare {
	return []gamelogic.PieceSquare{
		10, 11, 14, 15,
		18, 19, 22, 23,
	}
}

func getEdgeSquares() []gamelogic.PieceSquare {
	return []gamelogic.PieceSquare{
		1, 2, 3, 4,
		5, 13, 21, 29,
		12, 20, 28,
		30, 31, 32,
	}
}

func getBlackBackrank() []gamelogic.PieceSquare {
	return []gamelogic.PieceSquare{
		5, 6, 7, 8,
	}
}

func getRedBlackrank() []gamelogic.PieceSquare {
	return []gamelogic.PieceSquare{
		25, 26, 27, 28,
	}
}

type BotMinimaxConfig struct {
	Name          string
	ManWeight     float32
	KingWeight    float32
	CenterBonus   float32
	EdgeBonus     float32
	BackrankBonus float32
}

func GetDefaultBotMinimaxConfig() BotMinimaxConfig {
	return BotMinimaxConfig{
		Name:          "default",
		ManWeight:     1.0,
		KingWeight:    2.0,
		CenterBonus:   0.25,
		EdgeBonus:     0.25,
		BackrankBonus: 0.5,
	}
}

type BotMinimax struct {
	depth  uint
	config BotMinimaxConfig
}

func InitBotMinimax(depth uint, configParam ...BotMinimaxConfig) BotMinimax {
	var config BotMinimaxConfig
	if len(configParam) > 0 {
		config = configParam[0]
	} else {
		config = GetDefaultBotMinimaxConfig()
	}
	return BotMinimax{depth, config}
}

func (bot BotMinimax) Name() string {
	return fmt.Sprintf("minimax-%s-%d", bot.config.Name, bot.depth)
}

func (bot BotMinimax) Move(board gamelogic.Board) (gamelogic.BoardMove, bool) {
	// TODO: complete this function
	panic("Not implemented!")
}

func (bot BotMinimax) evaluate(board gamelogic.Board, player gamelogic.PieceColor) float32 {
	weight := float32(0.0)
	for _, piece := range board.Pieces() {
		coefficient := float32(1.0)
		if piece.Color != player {
			coefficient = -1.0
		}
		local := float32(0.0)
		switch piece.Type {
		case gamelogic.MAN:
			local += bot.config.ManWeight
		case gamelogic.KING:
			local += bot.config.KingWeight
		}
		for _, centerSquare := range getCenterSquares() {
			if piece.Square == centerSquare {
				local += bot.config.CenterBonus
				break
			}
		}
		for _, edgeSquare := range getEdgeSquares() {
			if piece.Square == edgeSquare {
				local += bot.config.EdgeBonus
				break
			}
		}
		var backrank []gamelogic.PieceSquare
		if piece.Color == gamelogic.BLACK {
			backrank = getBlackBackrank()
		} else {
			backrank = getRedBlackrank()
		}
		for _, backrankSquare := range backrank {
			if piece.Square == backrankSquare {
				local += bot.config.BackrankBonus
				break
			}
		}
		weight += local * coefficient
	}
	return weight
}
