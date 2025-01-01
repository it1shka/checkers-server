package multiplayer

import "it1shka.com/checkers-server/gamelogic"

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
