package multiplayer

type matchmaking struct {
	// TODO:
}

func newMatchmaking() *matchmaking {
	return &matchmaking{
		// TODO:
	}
}

func (m *matchmaking) handleAsync(player *player) {
	go m.handleJoin(player)
	go m.handleLeave(player)
	go m.handleMove(player)
}

func (m *matchmaking) handleJoin(player *player) {
	for range player.joinChannel {
		// TODO:
	}
}

func (m *matchmaking) handleLeave(player *player) {
	for range player.leaveChannel {
		// TODO:
	}
}

func (m *matchmaking) handleMove(player *player) {
	for move := range player.movesChannel {
		// TODO:
	}
}
