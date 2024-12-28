package multiplayer

import "time"

// Constants

const (
	MSG_IN_JOIN  = "join"
	MSG_IN_LEAVE = "leave"
	MSG_IN_MOVE  = "move"
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

// Implementation

type Multiplayer struct {
	players *SafeDict[string, IPlayer]
	hooks   *SafeDict[string, func(update IUpdate)]
	queue   *SafeSet[string]
}

func NewMultiplayer() *Multiplayer {
	return &Multiplayer{
		players: NewSafeDict[string, IPlayer](),
		hooks:   NewSafeDict[string, func(update IUpdate)](),
		queue:   NewSafeSet[string](),
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
	m.queue.Remove(id)
	m.hooks.Delete(id)
	m.players.Delete(id)
	// TODO: Cleanup queue and possibly running session
}

func (m *Multiplayer) PushUpdate(id string, update IUpdate) {
	switch update.Type() {
	case MSG_IN_JOIN:
		// TODO: if already in session, reject
		m.queue.Add(id)
	case MSG_IN_LEAVE:
		m.queue.Remove(id)
	case MSG_IN_MOVE:

	}
}

func (m *Multiplayer) startHumanSession(idA, idB string) {
	// TODO:
}

func (m *Multiplayer) startBotSession(id string) {
	// TODO:
}
