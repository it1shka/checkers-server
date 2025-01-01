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

func (m *matchmaking) startMatchmakingAsync(period time.Duration) {
	go func() {
		for range time.Tick(period) {
			m.queue.WithLock(func(_queue map[*player]bool) {
				defer clear(_queue)
				players := utils.Keys(_queue)
				for i := 0; i < len(players); i += 2 {
					if i == len(players)-1 {
						m.startBotGameAsync(players[i])
					} else {
						m.startHumanGameAsync(players[i], players[i+1])
					}
				}
			})
		}
	}()
}

func (m *matchmaking) startHumanGameAsync(playerA, playerB *player) {
	game := newGame(playerA, playerB)
	m.games.Put(playerA, game)
	m.games.Put(playerB, game)
	go func() {
		defer func() {
			m.games.Delete(game.playerBlack)
			m.games.Delete(game.playerRed)
		}()
		<-game.done
	}()
	game.startAsync()
}

func (m *matchmaking) startBotGameAsync(player *player) {
	// TODO:
	println("UNIMPLEMENTED")
}

func (m *matchmaking) handlePlayerAsync(player *player) {
	go func() {
		defer m.cleanupPlayer(player)
		<-player.done
	}()

	go func() {
	ListenJoin:
		for {
			select {
			case <-player.done:
				break ListenJoin
			case <-player.joinChannel:
				if !m.games.HasKey(player) {
					m.queue.Add(player)
				}
			}
		}
	}()

	go func() {
	ListenMove:
		for {
			select {
			case <-player.done:
				break ListenMove
			case move := <-player.movesChannel:
				m.games.IfExists(player, func(activeGame *game) {
					activeGame.moves <- authoredMove{
						author: player,
						move:   move,
					}
				})
			}
		}
	}()

	go func() {
	ListenLeave:
		for {
			select {
			case <-player.done:
				break ListenLeave
			case <-player.leaveChannel:
				m.cleanupPlayer(player)
			}
		}
	}()
}

func (m *matchmaking) cleanupPlayer(player *player) {
	m.games.IfExists(player, func(activeGame *game) {
		activeGame.retreat(player)
	})
	m.games.Delete(player)
	m.queue.Delete(player)
}
