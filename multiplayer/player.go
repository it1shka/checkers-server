package multiplayer

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"it1shka.com/checkers-server/gamelogic"
)

type player struct {
	conn         *websocket.Conn
	info         playerInfo
	id           uuid.UUID
	joinChannel  chan bool
	leaveChannel chan bool
	movesChannel chan incomingMove
	sendChannel  chan outcomingMessage
	done         chan bool
}

func newPlayer(conn *websocket.Conn, info playerInfo) (*player, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	player := &player{
		conn:         conn,
		info:         info,
		id:           id,
		joinChannel:  make(chan bool),
		leaveChannel: make(chan bool),
		movesChannel: make(chan incomingMove),
		sendChannel:  make(chan outcomingMessage),
		done:         make(chan bool),
	}
	return player, nil
}

func (p *player) startAsync() {
	go p.receive()
	go p.listen()
}

func (p *player) receive() {
	for msg := range p.sendChannel {
		if err := p.conn.WriteJSON(msg); err != nil {
			log.Println(err)
		}
	}
}

func (p *player) listen() {
	defer p.stop()
Listening:
	for {
		var message incomingMessage
		if err := p.conn.ReadJSON(&message); err != nil {
			if !websocket.IsCloseError(err) && websocket.IsUnexpectedCloseError(err) {
				log.Println(err)
			}
			break Listening
		}
		switch message.Type {
		case incMsgJoin:
			p.joinChannel <- true
		case incMsgLeave:
			p.leaveChannel <- true
		case incMsgMove:
			var moveMessage incomingMoveRaw
			if err := json.Unmarshal(message.Payload, &moveMessage); err != nil {
				continue Listening
			}
			from := gamelogic.PieceSquare(moveMessage.From)
			to := gamelogic.PieceSquare(moveMessage.To)
			if !from.IsValid() || !to.IsValid() {
				continue Listening
			}
			p.movesChannel <- incomingMove{from, to}
		}
	}
}

func (p *player) stop() {
	p.conn.Close()
	close(p.sendChannel)
	close(p.movesChannel)
	close(p.leaveChannel)
	close(p.joinChannel)
	close(p.done)
}
