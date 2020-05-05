package game

import (
	"errors"

	"github.com/wjase/cheads/pkg/websocket"
)

// WebsocketGameHub which uses websockets for player comms
type WebsocketGameHub struct {
	Pool    *websocket.Pool
	Players map[string]Player // mapped by id
	State   map[string]interface{}
}

// NewGameHub constructor
func NewGameHub(pool *websocket.Pool) *WebsocketGameHub {
	hub := WebsocketGameHub{Pool: pool, Players: map[string]Player{}, State: map[string]interface{}{}}
	return &hub
}

// AddPlayer adds a new player to the game
func (wg *WebsocketGameHub) AddPlayer(player Player) {
	wg.Players[player.ID()] = player
}

// GetPlayer by ID
func (wg *WebsocketGameHub) GetPlayer(ID string) (Player, error) {
	if player, ok := wg.Players[ID]; ok {
		return player, nil
	}
	return nil, errors.New("No player found")

}

// Start begin the game
func (wg *WebsocketGameHub) Start() error {
	panic("not implemented") // TODO: Implement
}

// Reset begin the game again
func (wg *WebsocketGameHub) Reset() error {
	panic("not implemented") // TODO: Implement
}

//GetState get the state
func (wg *WebsocketGameHub) GetState() map[string]interface{} {
	panic("not implemented") // TODO: Implement
}
