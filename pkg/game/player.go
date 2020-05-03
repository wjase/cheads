package game

//Player generic game player
type Player interface {
	ID() string
	Name() string
	Team() string
	SendMessage(message string) error
	OnMessage(message string) error
	GetState() map[string]interface{}
	SetState(map[string]interface{})
}
