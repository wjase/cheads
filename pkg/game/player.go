package game

//go:generate mockgen -destination=../mock_game/MockPlayer.go	-package=mock_game . Player

//Player generic game player
type Player interface {
	// The unique id of this player
	ID() string
	// The human readable name of this player
	Name() string
	// The human readable team name of this player
	Team() string
	// sends a message to this Player
	SendEvent(event Event) error

	//
	GetState(key string) interface{}
	PutState(key string, state interface{})
	GetStates() map[string]interface{}
}
