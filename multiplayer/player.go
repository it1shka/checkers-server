package multiplayer

import (
	"github.com/gorilla/websocket"
)

type Player struct {
	Conn     *websocket.Conn
	Nickname string `schema:"nickname"`
	Rating   uint   `schema:"rating"`
	Region   string `schema:"region"`
}
