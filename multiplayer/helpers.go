package multiplayer

import (
	"math/rand"
	"strconv"
	"time"

	"it1shka.com/checkers-server/gamelogic"
)

func getOutMsgEnemy(enemy *player) outcomingMessage {
	return outcomingMessage{
		Type:    outMsgEnemy,
		Payload: enemy.info,
	}
}

func getOutMsgColor(color string) outcomingMessage {
	return outcomingMessage{
		Type:    outMsgColor,
		Payload: color,
	}
}

func getOutMsgBoard(board gamelogic.Board) outcomingMessage {
	type jsonPiece struct {
		Color  string `json:"color"`
		Square int    `json:"square"`
		Type   string `json:"type"`
	}

	jsonPieces := make([]jsonPiece, len(board.UnsafePieces()))
	for index, piece := range board.UnsafePieces() {
		jsonPieces[index] = jsonPiece{
			Color:  piece.Color.String(),
			Square: int(piece.Square),
			Type:   piece.Type.String(),
		}
	}

	return outcomingMessage{
		Type: outMsgBoard,
		Payload: map[string]any{
			"turn":   board.Turn().String(),
			"pieces": jsonPieces,
		},
	}
}

func getOutMsgStatus(status gamelogic.GameStatus) outcomingMessage {
	var statusString string

	switch status {
	case gamelogic.ACTIVE:
		statusString = "active"
	case gamelogic.DRAW:
		statusString = "draw"
	case gamelogic.BLACK_WIN:
		statusString = "black"
	case gamelogic.RED_WIN:
		statusString = "red"
	}

	return outcomingMessage{
		Type:    outMsgStatus,
		Payload: statusString,
	}
}

func getOutMsgTime(tag string, time int32) outcomingMessage {
	return outcomingMessage{
		Type: outMsgTime,
		Payload: map[string]any{
			"player": tag,
			"time":   time,
		},
	}
}

func getOutMsgQueueJoined() outcomingMessage {
	return outcomingMessage{
		Type: outMsgQueueJoined,
	}
}

func getOutMsgQueueLeft() outcomingMessage {
	return outcomingMessage{
		Type: outMsgQueueLeft,
	}
}

// TODO: these names are pretty weird
var pseudoFirstNames = []string{
	"Astro",
	"Frank",
	"Hugo",
	"Big",
	"Profi",
	"Sn1per",
	"pr0fessional",
}

var pseudoSecondNames = []string{
	"_profi",
	"__tank",
	"Master",
	"TopPlayer",
	"Messi",
	"Krash",
	"@storm",
}

var pseudoCountries = []string{
	"Belarus",
	"Poland",
	"Germany",
	"US",
	"France",
	"Australia",
	"Hungary",
}

func getPseudoPlayerInfo() playerInfo {
	firstPart := pseudoFirstNames[rand.Intn(len(pseudoFirstNames))]
	secondPart := pseudoSecondNames[rand.Intn(len(pseudoSecondNames))]
	postfix := rand.Intn(900) + 100
	name := firstPart + secondPart + strconv.Itoa(postfix)
	rating := rand.Intn(101)
	country := pseudoCountries[rand.Intn(len(pseudoCountries))]
	return playerInfo{
		Nickname: name,
		Rating:   uint(rating),
		Region:   country,
	}
}

func getOutMsgPseudoEnemy(info playerInfo) outcomingMessage {
	return outcomingMessage{
		Type:    outMsgEnemy,
		Payload: info,
	}
}

func randomBool() bool {
	seed := time.Now().UnixNano()
	gen := rand.New(rand.NewSource(seed))
	guess := gen.Intn(2)
	return guess == 0
}
