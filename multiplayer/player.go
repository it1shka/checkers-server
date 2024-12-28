package multiplayer

import "it1shka.com/checkers-server/gamelogic"

// Abstract player = human or bot

type PlayerInfo interface {
	Name() string
	Ranking() uint
	Region() string
}

type Player interface {
	ID() string
	Info() PlayerInfo
	Start(handler func(move gamelogic.BoardMove))
	NotifyAboutEnemy(info PlayerInfo)
	NotifyAboutUpdate(update Update)
	Stop()
}

// Concrete player implementations

// TODO:

type PlayerHuman struct {
	info        PlayerInfo
	moves       <-chan gamelogic.BoardMove
	enemyUpdate chan<- PlayerInfo
	boardUpdate chan<- Update
}

func NewPlayerHuman(info PlayerInfo) *PlayerHuman {
	return &PlayerHuman{}
}
