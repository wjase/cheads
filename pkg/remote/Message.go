package remote

// Message Struct for wrapping messages
type Message struct {
	FromID string `json:"from"`
	ToID   string `json:"to"`
	Body   string `json:"body"`
}
