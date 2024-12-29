package multiplayer

import (
	"time"

	"it1shka.com/checkers-server/utils"
)

type matchmaking struct {
	queue *utils.SafeSet[*player]
	games *utils.SafeDict[*player, *game]
}

func newMatchmaking() *matchmaking {
	return &matchmaking{
		queue: utils.NewSafeSet[*player](),
		games: utils.NewSafeDict[*player, *game](),
	}
}

func (m *matchmaking) handleQueue(period time.Duration) {
	for range time.Tick(period) {
		m.queue.WithLock(func(_queue map[*player]bool) {
			defer clear(_queue)
			players := utils.Keys(_queue)
			for i := 0; i < len(players); i += 2 {
				if i == len(players)-1 {
					m.startBotGame(players[i])
				} else {
					m.startHumanGame(players[i], players[i+1])
				}
			}
		})
	}
}

func (m *matchmaking) startHumanGame(playerA, playerB *player) {
	game := newGame(playerA, playerB)
	m.games.Put(playerA, game)
	m.games.Put(playerB, game)
	game.startAsync()
	go m.cleanupGame(game)
}

func (m *matchmaking) startBotGame(player *player) {
	// TODO:
	println("UNIMPLEMENTED")
}

func (m *matchmaking) cleanupGame(game *game) {
	defer func() {
		m.games.Delete(game.playerBlack)
		m.games.Delete(game.playerRed)
	}()
	<-game.done
}

func (m *matchmaking) handleAsync(player *player) {
	go m.handleCleanup(player)
	go m.handleJoin(player)
	go m.handleLeave(player)
}

func (m *matchmaking) handleJoin(player *player) {
	for range player.joinChannel {
		if !m.games.HasKey(player) {
			m.queue.Add(player)
		}
	}
}

func (m *matchmaking) handleLeave(player *player) {
	for range player.leaveChannel {
		m.cleanupPlayer(player)
	}
}

func (m *matchmaking) handleCleanup(player *player) {
	defer m.cleanupPlayer(player)
	<-player.done
}

func (m *matchmaking) cleanupPlayer(player *player) {
	m.games.IfExists(player, func(activeGame *game) {
		activeGame.retreat(player)
	})
	m.games.Delete(player)
	m.queue.Delete(player)
}
