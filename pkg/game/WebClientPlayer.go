package game

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/wjase/cheads/pkg/remote"
)

// WebClientPlayer a Player backed by a ws Client connection
type WebClientPlayer struct {
	Client     remote.Client
	State      map[string]interface{}
	PlayerName string
	PlayerTeam string
}

// NewWebClientPlayer creates a running Client
func NewWebClientPlayer(client remote.Client) *WebClientPlayer {
	player := WebClientPlayer{Client: client, State: map[string]interface{}{}}
	client.Subscribe(&player)
	return &player
}

//SendEvent sends a message to this Player's client
func (wp *WebClientPlayer) SendEvent(event Event) error {
	b, err := json.Marshal(event)
	if err != nil {
		log.Println("error:", err)
		return err
	}
	msg := remote.Message{Body: string(b), FromID: event.FromID, ToID: wp.ID()}
	wp.Client.Send(msg)
	return nil
}

// ID The unique id of this player
func (wp *WebClientPlayer) ID() string {
	return wp.Client.ID()
}

// Name The human readable name of this player
func (wp *WebClientPlayer) Name() string {
	return wp.PlayerName
}

// Team The human readable team name of this player
func (wp *WebClientPlayer) Team() string {
	return wp.PlayerTeam
}

// OnMessage implementation of Client fn
func (wp *WebClientPlayer) OnMessage(message remote.Message) error {
	evt := Event{}
	err := json.NewDecoder(strings.NewReader(message.Body)).Decode(&evt)
	if err != nil {
		log.Printf("Unable to parse message %s", message)
		return err
	}

	// publish the event to a listener
	return nil
}

// GetState return the named state item
func (wp *WebClientPlayer) GetState(key string) interface{} {
	return wp.State[key]
}

// GetStates return all the states
func (wp *WebClientPlayer) GetStates() map[string]interface{} {
	return wp.State
}

// SetStates replace the state map for the player
func (wp *WebClientPlayer) SetStates(toSet map[string]interface{}) {
	wp.State = toSet
}

// PutState TODO COMMENT HERE
func (wp *WebClientPlayer) PutState(key string, state interface{}) {
	wp.State[key] = state
}
