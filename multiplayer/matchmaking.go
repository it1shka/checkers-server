package multiplayer

import (
	"time"

	"it1shka.com/checkers-server/utils"
)

type abstractGame interface {
	retreat(player *player)
	pushMove(move authoredMove)
}

type matchmaking struct {
	queue *utils.SafeSet[*player]
	games *utils.SafeDict[*player, abstractGame]
}

func newMatchmaking() *matchmaking {
	return &matchmaking{
		queue: utils.NewSafeSet[*player](),
		games: utils.NewSafeDict[*player, abstractGame](),
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
	pseudogame := newPseudogame(player)
	m.games.Put(player, pseudogame)
	go func() {
		defer m.games.Delete(player)
		<-pseudogame.done
	}()
	pseudogame.startAsync()
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
				m.games.IfExists(player, func(activeGame abstractGame) {
					activeGame.retreat(player)
				})
				m.games.Delete(player)
				m.queue.Add(player)
				player.sendMessage(getOutMsgQueueJoined())
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
				m.games.IfExists(player, func(activeGame abstractGame) {
					move := authoredMove{
						author: player,
						move:   move,
					}
					activeGame.pushMove(move)
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
				if m.queue.Has(player) {
					player.sendMessage(getOutMsgQueueLeft())
				}
				m.cleanupPlayer(player)
			}
		}
	}()
}

func (m *matchmaking) cleanupPlayer(player *player) {
	m.games.IfExists(player, func(activeGame abstractGame) {
		activeGame.retreat(player)
	})
	m.games.Delete(player)
	m.queue.Delete(player)
}
