package multiplayer

import (
	"encoding/json"
	"time"

	"it1shka.com/checkers-server/gamelogic"
)

// Constants

const (
	MSG_IN_JOIN  = "join"
	MSG_IN_LEAVE = "leave"
	MSG_IN_MOVE  = "move"
)

const (
	MSG_OUT_ENEMY  = "enemy"
	MSG_OUT_UPDATE = "update"
)

// Helper records

type IPlayer interface {
	ID() string
	Nickname() string
	Rating() uint
	Region() string
}

type IUpdate interface {
	Type() string
	Payload() any
}

type update struct {
	Type_    string `json:"type"`
	Payload_ any    `json:"payload"`
}

func (u update) Type() string {
	return u.Type_
}
func (u update) Payload() any {
	return u.Payload_
}

// Implementation

type Multiplayer struct {
	players  *SafeDict[string, IPlayer]
	hooks    *SafeDict[string, func(update IUpdate)]
	queue    *SafeSet[string]
	sessions *SafeDict[string, *Session]
}

func NewMultiplayer() *Multiplayer {
	return &Multiplayer{
		players:  NewSafeDict[string, IPlayer](),
		hooks:    NewSafeDict[string, func(update IUpdate)](),
		queue:    NewSafeSet[string](),
		sessions: NewSafeDict[string, *Session](),
	}
}

func (m *Multiplayer) StartMatchmaking(period time.Duration) {
	go func() {
		for range time.Tick(period) {
			ids := m.queue.Eject()
			for i := 0; i < len(ids); i += 2 {
				if i == len(ids)-1 {
					m.startBotSession(ids[i])
				} else {
					m.startHumanSession(ids[i], ids[i+1])
				}
			}
		}
	}()
}

func (m *Multiplayer) RegisterPlayer(player IPlayer, hook func(update IUpdate)) {
	m.players.Put(player.ID(), player)
	m.hooks.Put(player.ID(), hook)
}

func (m *Multiplayer) UnregisterPlayer(id string) {
	m.sessions.IfExists(id, func(session *Session) {
		session.SurrenderAndStop(id)
	})
	m.sessions.Delete(id)
	m.queue.Remove(id)
	m.hooks.Delete(id)
	m.players.Delete(id)
}

func (m *Multiplayer) PushUpdate(id string, update IUpdate) {
	switch update.Type() {
	case MSG_IN_JOIN:
		if !m.sessions.HasKey(id) {
			m.queue.Add(id)
		}
	case MSG_IN_LEAVE:
		m.sessions.IfExists(id, func(session *Session) {
			session.SurrenderAndStop(id)
		})
		m.sessions.Delete(id)
		m.queue.Remove(id)
	case MSG_IN_MOVE:
		// TODO: this is a shitty workaround,
		// TODO: remove in the future
		// TODO: maybe add MapToStruct function
		// TODO: that will use reflection
		data, err := json.Marshal(update.Payload)
		if err != nil {
			return
		}
		var move struct {
			From uint `json:"from"`
			To   uint `json:"to"`
		}
		if err := json.Unmarshal(data, &move); err != nil {
			return
		}
		m.sessions.IfExists(id, func(session *Session) {
			session.Move(SessionMove{
				Player: id,
				From:   move.From,
				To:     move.To,
			})
		})
	}
}

func (m *Multiplayer) startHumanSession(idA, idB string) {
	m.hooks.IfExists(idA, func(hook func(update IUpdate)) {
		playerInfo := m.players.GetOrEmpty(idB)
		hook(update{
			Type_: MSG_OUT_ENEMY,
			Payload_: map[string]any{
				"nickname": playerInfo.Nickname(),
				"rating":   playerInfo.Rating(),
				"region":   playerInfo.Region(),
			},
		})
	})
	m.hooks.IfExists(idB, func(hook func(update IUpdate)) {
		playerInfo := m.players.GetOrEmpty(idA)
		hook(update{
			Type_: MSG_OUT_ENEMY,
			Payload_: map[string]any{
				"nickname": playerInfo.Nickname(),
				"rating":   playerInfo.Rating(),
				"region":   playerInfo.Region(),
			},
		})
	})

	session := NewSession(idA, idB, func(upd SessionUpdate) {
		var statusString string
		switch upd.Status {
		case gamelogic.ACTIVE:
			statusString = "active"
		case gamelogic.DRAW:
			statusString = "draw"
		case gamelogic.BLACK_WIN:
			statusString = "black"
		case gamelogic.RED_WIN:
			statusString = "red"
		}

		rawPieces := upd.Board.UnsafePieces()
		type piece struct {
			Square int    `json:"square"`
			Color  string `json:"color"`
			Type   string `json:"type"`
		}
		pieces := make([]piece, len(rawPieces))
		for index, rawPiece := range rawPieces {
			var colorString string
			switch rawPiece.Color {
			case gamelogic.RED:
				colorString = "red"
			default:
				colorString = "black"
			}
			var typeString string
			switch rawPiece.Type {
			case gamelogic.MAN:
				typeString = "man"
			default:
				typeString = "king"
			}
			pieces[index] = piece{
				Square: int(rawPiece.Square),
				Color:  colorString,
				Type:   typeString,
			}
		}

		payload := map[string]any{
			"status": statusString,
			"board":  pieces,
		}
		msgUpdate := update{
			Type_:    MSG_OUT_UPDATE,
			Payload_: payload,
		}

		m.hooks.IfExists(idA, func(hook func(update IUpdate)) {
			hook(msgUpdate)
		})
		m.hooks.IfExists(idB, func(hook func(update IUpdate)) {
			hook(msgUpdate)
		})
	})
	m.sessions.Put(idA, session)
	m.sessions.Put(idB, session)
	session.Start()
}

func (m *Multiplayer) startBotSession(id string) {
	// TODO:
}
