package celebgame

import "github.com/wjase/cheads/pkg/websocket"

// Engine implements the rules and state changes for Celeb Game
type Engine struct {
	Pool *websocket.Pool
}

// NewGame crates a nw game instancee
func NewGame(pool *websocket.Pool) *Engine {
	engine := Engine{Pool: pool}
	return &engine
}
