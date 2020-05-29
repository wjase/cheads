package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/wjase/cheads/pkg/remote"
)

// WsHub which uses websockets for player comms
type WsHub struct {
	Pool       remote.ClientPool
	Players    map[string]Player // mapped by id
	State      map[string]interface{}
	Events     chan *Event
	Listeners  []interface{}
	GameEngine Engine
}

// NewGameHub constructor
func NewGameHub(pool remote.ClientPool) *WsHub {
	hub := WsHub{Pool: pool, Players: map[string]Player{}, State: map[string]interface{}{}}
	pool.SubscribeJoined(hub.ClientJoined)
	return &hub
}

// AddGame adds and starts the game engine
func (wg *WsHub) AddGame(engine Engine) {
	wg.GameEngine = engine
	go func() {
		err := wg.GameEngine.Run()
		fmt.Printf("ERROR:%v\n", err)
	}()
}

// ClientJoined When a client joins the pool
func (wg *WsHub) ClientJoined(client remote.Client) {
	player := WebClientPlayer{Client: client, PlayerName: client.ID(), PlayerTeam: client.ID(), State: map[string]interface{}{}}
	wg.AddPlayer(&player)
}

// AddPlayer adds a new player to the game
func (wg *WsHub) AddPlayer(player Player) {
	wg.Players[player.ID()] = player
	wg.Publish(&Event{FromID: "game", PlayerJoined: &PlayerJoinedEvent{Player: player}})
}

// GetPlayer by ID
func (wg *WsHub) GetPlayer(ID string) (Player, error) {
	if player, ok := wg.Players[ID]; ok {
		return player, nil
	}
	return nil, errors.New("No player found")
}

// GetPlayers gets the list of players
func (wg *WsHub) GetPlayers() map[string]Player {
	return wg.Players
}

// Start begin the game
func (wg *WsHub) Start() error {
	return wg.GameEngine.Run()
}

// Reset begin the game again
func (wg *WsHub) Reset() error {
	return wg.GameEngine.Reset()
}

//GetState get the state
func (wg *WsHub) GetState() map[string]interface{} {
	return wg.State
}

//Publish send message
func (wg *WsHub) Publish(event *Event) {
	byt, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error when publishing: %v", err)
	}
	wg.Pool.Dispatch(remote.Message{Body: string(byt)})
}

// Subscribe should subscribe to messages
func (wg *WsHub) Subscribe(listener interface{}) {
	wg.Listeners = append(wg.Listeners, listener)
}
