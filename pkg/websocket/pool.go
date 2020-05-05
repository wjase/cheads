package websocket

import (
	"fmt"
	"net/http"

	"github.com/wjase/cheads/pkg/auth"
)

// ClientJoined event callback when Client Joined
type ClientJoined func(*Client)

// ClientLeft event callback when Client Left
type ClientLeft func(*Client)

// Pool mmanages creation of Clients and messaging between them
type Pool struct {
	Register          chan *Client
	Unregister        chan *Client
	Clients           map[string]*Client
	Broadcast         chan Message
	ClientJoinedGroup []ClientJoined
	ClientLeftGroup   []ClientLeft
}

//NewPool create new Pool
func NewPool() *Pool {
	return &Pool{
		Register:          make(chan *Client),
		Unregister:        make(chan *Client),
		Clients:           map[string]*Client{},
		Broadcast:         make(chan Message),
		ClientLeftGroup:   []ClientLeft{},
		ClientJoinedGroup: []ClientJoined{},
	}
}

//AddClient add a client to the list
func (pool *Pool) AddClient(client *Client) {
	pool.Clients[client.ID] = client
}

func getClaimsFromRequestContext(r *http.Request) *auth.Claims {
	val := r.Context().Value(auth.ContextClaimsKey)
	claims, ok := val.(*auth.Claims)
	if ok {
		return claims
	}
	fmt.Printf("Unable to convert %v into a Claims object", val)
	return nil
}

// ServeWsCreateClient http handler to upgrade a request to websocket
// and create a new Client, then run it's event loop until the connection dies
func (pool *Pool) ServeWsCreateClient() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := getClaimsFromRequestContext(r)

		if claims == nil {
			fmt.Printf("Missing expected claims context on request")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		conn, err := Upgrade(w, r)
		if err != nil {
			fmt.Fprintf(w, "%+v\n", err)
			w.WriteHeader(http.StatusBadRequest)
		}

		client := &Client{
			ID:   claims.Username,
			Conn: conn,
			Pool: pool,
		}

		pool.Register <- client
		client.ReadLoop()
	}
}

//Send send a message to all Clients
func (pool *Pool) Send(message Message) {
	for _, client := range pool.Clients {
		if err := client.Conn.WriteJSON(message); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.AddClient(client)
			pool.Send(Message{Type: 1, Body: "New User Joined..."})
			break
		case client := <-pool.Unregister:
			delete(pool.Clients, client.ID)
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			pool.Send(Message{Type: 1, Body: "User Disconnected..."})
			break
		case message := <-pool.Broadcast:
			fmt.Println("Sending message to all clients in Pool")
			pool.Send(message)
		}
	}
}
