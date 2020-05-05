package game

// PlayerAdded fired when a player has joined the game
type PlayerAdded func(player Player)

// PlayerUpdated fired when a player leaves the game
type PlayerUpdated func(player Player)

// Fired when a game starts
type Started func()

// Fired When a game state is updated
type Updated func()

// Hub - a shared game space
// A hub has players who each have their own state (eg piece position , cards etc)
// The hub may have it's own state: eg card decks discard piles etc (current player/team etc)
type Hub interface {
	AddPlayer(player Player)
	GetPlayer(ID string) (Player, error)
	Start() error
	Reset() error
	GetState() map[string]interface{}
	GetStateFor(id string) interface{}
	PutStateFor(id string, state interface{})
	RegisterPlayerAdded(fn PlayerAdded)
	RegisterPlayerUpdated(fn PlayerUpdated)
	RegisterStarted(fn Started)
	RegisterUpdated(fn Started)
}
