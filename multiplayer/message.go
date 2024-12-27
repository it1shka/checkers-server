package multiplayer

import "encoding/json"

type ClientMessage struct {
  Type string `json:"type"`
  Payload json.RawMessage `json:"payload"`
}

type ServerMessage struct {
  Type string `json:"type"`
  Payload any `json:"payload"`
}
