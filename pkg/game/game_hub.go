package game

import "fmt"

//go:generate mockgen -destination=../mock_game/MockPlayerHub.go	-package=mock_game . PlayerHub

// PlayerHub for game communication
type PlayerHub interface {
	// Send event to all Players
	Broadcast(msg *Event)
	// Send Event to single Player
	SendEvent(toID string, evt Event)

	Subscribe(listener interface{})

	// Gets a specific player
	GetPlayer(ID string) (Player, error)
	GetPlayers() map[string]Player
	// // PlayerJoined called when a new player joins the game
	// PlayerJoined(player *Player)

	// // PlayerLeft called when players leave
	// PlayerLeft(player *Player)
}

// Engine Generic game logic implemented here
type Engine interface {
	// Runs the GameLoop
	Run() error

	// Reset reset the game state
	Reset() error

	// OnEvent when a game event occurs
	OnEvent(evt *Event)
}

// PlayerLeftListener implement this to respond to left events
type PlayerLeftListener interface {
	PlayerLeft(player Player)
}

// Event generic event envelope
// The message will contain one non-null payload
type Event struct {
	// The source of the event
	FromID        string
	PlayerJoined  *PlayerJoinedEvent  `json:"player_joined"`
	PlayerLeft    *PlayerLeftEvent    `json:"player_left"`
	PlayerUpdated *PlayerUpdatedEvent `json:"player_updated"`
	Updated       *UpdatedEvent       `json:"game_updated"`
	TextMessage   *TextMessageEvent   `json:"text_message"`
}

// TriggerOn If the receiver can handle the event then fire it
// if not ignore it
// In order to handle an event, simply implment the right handler method
// on the receiver
func (e *Event) TriggerOn(receiver interface{}) {
	switch {
	case e.PlayerJoined != nil:
		if listener, ok := receiver.(PlayerJoinedListener); ok {
			listener.PlayerJoined(e.PlayerJoined.Player)
		}
	case e.PlayerLeft != nil:
		if listener, ok := receiver.(PlayerLeftListener); ok {
			listener.PlayerLeft(e.PlayerLeft.Player)
		}
	case e.PlayerUpdated != nil:
		if listener, ok := receiver.(PlayerUpdatedListener); ok {
			listener.PlayerUpdated(e.FromID, e.PlayerUpdated.Player)
		}
	case e.Updated != nil:
		if listener, ok := receiver.(UpdatedListener); ok {
			listener.Updated(e.FromID, e.Updated.StateID, e.Updated.State)
		}
	case e.TextMessage != nil:
		if listener, ok := receiver.(TextMessageListener); ok {
			listener.TextMessage(e.FromID, e.TextMessage.ToID, e.TextMessage.MessageText)
		}

	default:
		fmt.Printf("Unnknown event type")
	}

}

// TextMessageListener implement this to respond to joined events
type TextMessageListener interface {
	TextMessage(FromID string, ToID string, MessageText string)
}

// TextMessageEvent Used for player chat
type TextMessageEvent struct {
	ToID        string `json:"to_id"`
	MessageText string `json:"message_text"`
}

// PlayerJoinedListener implement this to respond to joined events
type PlayerJoinedListener interface {
	PlayerJoined(player Player)
}

// PlayerJoinedEvent - when a player is added
type PlayerJoinedEvent struct {
	Player Player `json:"player"`
}

// PlayerLeftEvent - when a player leaves
type PlayerLeftEvent struct {
	Player Player `json:"player"`
}

// PlayerUpdatedListener implement this to respond to joined events
type PlayerUpdatedListener interface {
	PlayerUpdated(FromID string, player Player)
}

// PlayerUpdatedEvent - when a player is updated
type PlayerUpdatedEvent struct {
	Player  Player      `json:"player"`
	StateID string      `json:"state_id"`
	State   interface{} `json:"state"`
}

// UpdatedListener implement this to respond to joined events
type UpdatedListener interface {
	Updated(FromID string, StateID string, State interface{})
}

// UpdatedEvent Fired When a game state is updated
type UpdatedEvent struct {
	StateID string      `json:"state_id"`
	State   interface{} `json:"state"`
}
