package multiplayer

import (
	"encoding/json"

	"it1shka.com/checkers-server/gamelogic"
)

type playerInfo struct {
	Nickname string `json:"nickname" schema:"nickname"`
	Rating   uint   `json:"rating" schema:"rating"`
	Region   string `json:"region" schema:"region"`
}

const (
	incMsgJoin  = "join"
	incMsgLeave = "leave"
	incMsgMove  = "move"
)

type incomingMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type incomingMoveRaw struct {
	From uint `json:"from"`
	To   uint `json:"to"`
}

type incomingMove struct {
	from gamelogic.PieceSquare
	to   gamelogic.PieceSquare
}

const (
	outMsgEnemy  = "enemy"
	outMsgColor  = "color"
	outMsgBoard  = "board"
	outMsgStatus = "status"
	outMsgTime   = "time"
)

type outcomingMessage struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

type authoredMove struct {
	author *player
	move   incomingMove
}
