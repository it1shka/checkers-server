package multiplayer

type Matchmaking[P any] struct {
	players chan P
	action  func(playerA, playerB P)
}

func NewMatchmaking[P any](action func(playerA, playerB P)) *Matchmaking[P] {
	return &Matchmaking[P]{
		players: make(chan P),
		action:  action,
	}
}

func (m *Matchmaking[P]) Enqueue(player P) {
	m.players <- player
}

func (m *Matchmaking[P]) Start() {
	var nullable P

	go func() {
		first := nullable
		active := false
		for player := range m.players {
			if !active {
				first = player
				active = true
			} else {
				m.action(first, player)
				first = nullable
				active = false
			}
		}
	}()
}

func (m *Matchmaking[P]) Stop() {
	close(m.players)
}
