package game

// Hub - a shared game space
// A hub has players who each have their own state (eg piece position , cards etc)
// The hub may have it's own state: eg card decks discard piles etc (current player/team etc)
type Hub interface {
	AddPlayer(player Player)
	GetPlayer(ID string) (Player, error)
	Start() error
	Reset() error
	GetState() map[string]interface{}
}
