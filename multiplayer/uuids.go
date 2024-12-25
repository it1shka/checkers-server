package multiplayer

import (
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type UUIDs struct {
	mutex   sync.RWMutex
	storage map[*websocket.Conn]uuid.UUID
	lookup  map[uuid.UUID]*websocket.Conn
}

func InitUUIDs() UUIDs {
	return UUIDs{
		mutex:   sync.RWMutex{},
		storage: make(map[*websocket.Conn]uuid.UUID),
		lookup:  make(map[uuid.UUID]*websocket.Conn),
	}
}

func (uuids *UUIDs) GenerateFor(conn *websocket.Conn) (uuid.UUID, bool) {
	newId := uuid.New()
	if newId.String() == "" {
		return newId, false
	}
	uuids.mutex.Lock()
	defer uuids.mutex.Unlock()
	uuids.storage[conn] = newId
	uuids.lookup[newId] = conn
	return newId, true
}

func (uuids *UUIDs) GetFor(conn *websocket.Conn) (uuid.UUID, bool) {
	uuids.mutex.RLock()
	defer uuids.mutex.RUnlock()
	result, ok := uuids.storage[conn]
	return result, ok
}

func (uuids *UUIDs) LookupFor(id uuid.UUID) (*websocket.Conn, bool) {
	uuids.mutex.RLock()
	defer uuids.mutex.RUnlock()
	result, ok := uuids.lookup[id]
	return result, ok
}

func (uuids *UUIDs) DeactivateFor(conn *websocket.Conn) {
	if id, ok := uuids.storage[conn]; ok {
		delete(uuids.lookup, id)
	}
	delete(uuids.storage, conn)
}
