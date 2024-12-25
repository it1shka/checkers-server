package multiplayer

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Player struct {
	Conn     *websocket.Conn
	ID       string `schema:"id"`
	Nickname string `schema:"nickname"`
	Rating   uint   `schema:"rating"`
}

type PlayerCollection struct {
	mutex   sync.RWMutex
	storage map[*websocket.Conn]Player
}

func NewPlayerCollection() *PlayerCollection {
	return &PlayerCollection{
		mutex:   sync.RWMutex{},
		storage: make(map[*websocket.Conn]Player),
	}
}

func (p *PlayerCollection) RegisterPlayer(player Player) bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, exists := p.storage[player.Conn]; exists {
		return false
	}
	p.storage[player.Conn] = player
	return true
}

func (p *PlayerCollection) DeletePlayerBy(conn *websocket.Conn) bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, exists := p.storage[conn]; !exists {
		return false
	}
	delete(p.storage, conn)
	return true
}

func (p *PlayerCollection) LookupPlayer(conn *websocket.Conn) (Player, bool) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	player, exists := p.storage[conn]
	return player, exists
}

func (p *PlayerCollection) GetSnapshot() []Player {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	output := make([]Player, len(p.storage))
	index := 0
	for _, player := range p.storage {
		output[index] = player
		index++
	}
	return output
}

func (p *PlayerCollection) Traverse(visitor func(Player)) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	for _, player := range p.storage {
		visitor(player)
	}
}
