package websocket

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/wjase/cheads/pkg/auth"
	"github.com/wjase/cheads/pkg/remote"
)

// Pool mmanages creation of Clients and messaging between them
type Pool struct {
	PoolID              string
	Register            chan remote.Client
	Unregister          chan remote.Client
	Clients             map[string]remote.Client
	Broadcast           chan remote.Message
	Messages            chan remote.Message
	ClientJoinedPoolFns []remote.ClientJoinedFn
	ClientLeftPoolFns   []remote.ClientLeftFn
}

// PoolSet Manages a set of pools
type PoolSet struct {
	Pools     map[string]remote.ClientPool
	PoolAdded remote.PoolAddedFn
}

func NePoolSet() *PoolSet {
	return &PoolSet{Pools: make(map[string]remote.ClientPool)}
}

var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

const roomCodeLen = 4

func (ps *PoolSet) nextRoomCode() string {
	for {
		b := make([]rune, roomCodeLen)
		for i := range b {
			b[i] = letters[rand.Intn(len(letters))]
		}
		code := string(b)
		// if not already used
		if _, ok := ps.Pools[code]; !ok {
			return code
		}
	}
}

//NewPool create new Pool
func (ps *PoolSet) NewPool() *Pool {
	pool := Pool{
		PoolID:     ps.nextRoomCode(),
		Register:   make(chan remote.Client),
		Unregister: make(chan remote.Client),
		Clients:    map[string]remote.Client{},
		// Broadcast:        make(chan Message),
		Messages:            make(chan remote.Message),
		ClientLeftPoolFns:   []remote.ClientLeftFn{ps.onClientLeft},
		ClientJoinedPoolFns: []remote.ClientJoinedFn{},
	}
	go pool.Start()
	return &pool
}

// ID from remote.ClientPool
func (pool *Pool) ID() string {
	return pool.PoolID
}

func (ps *PoolSet) onClientLeft(client remote.Client) {
	if client.Pool().Count() == 0 {
		ps.RemovePool(client.Pool())
	}
}

// RemovePool removes an empty pool from the set
func (ps *PoolSet) RemovePool(pool remote.ClientPool) {
	delete(ps.Pools, pool.ID())
}

//SubscribeJoined add a callback to be called when a client joins
func (pool *Pool) SubscribeJoined(joinedFn remote.ClientJoinedFn) {
	pool.ClientJoinedPoolFns = append(pool.ClientJoinedPoolFns, joinedFn)
}

//AddClient add a client to the list
func (pool *Pool) AddClient(client remote.Client) {
	pool.Clients[client.ID()] = client
	pool.BroadcastMsg(remote.Message{FromID: client.ID(), ToID: "ALL", Body: "client-joined"})
	for _, joinedFn := range pool.ClientJoinedPoolFns {
		joinedFn(client)
	}
}

// RemoveClient removes a client from the pool
func (pool *Pool) RemoveClient(client remote.Client) {
	delete(pool.Clients, client.ID())
	fmt.Println("Size of Connection Pool: ", len(pool.Clients))
	pool.BroadcastMsg(remote.Message{FromID: client.ID(), ToID: "ALL", Body: "client-left"})

	for _, leftFn := range pool.ClientLeftPoolFns {
		leftFn(client)
	}
}

func (pool *Pool) Dispatch(msg remote.Message) {
	pool.Messages <- msg
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

func (ps *PoolSet) poolForCode(code string) remote.ClientPool {
	if pool, ok := ps.Pools[code]; ok {
		return pool
	}
	pool := ps.NewPool()
	ps.Pools[pool.PoolID] = pool
	if ps.PoolAdded != nil {
		ps.PoolAdded(pool)
	}
	return pool
}

// ServeWsCreateClient http handler to upgrade a request to websocket
// and create a new Client, then run it's event loop until the connection dies
func (ps *PoolSet) ServeWsCreateClient() func(w http.ResponseWriter, r *http.Request) {
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

		// if joining exising game grab it and add client
		// if new game then create a new pool and give it a code
		code := ps.getRoomCodeForRequest(r)
		pool := ps.poolForCode(code)

		client := ps.createSocketPoolClient(conn, claims.Username, pool)
		client.ReadLoop()
	}
}

func (ps *PoolSet) getRoomCodeForRequest(r *http.Request) string {
	query := r.URL.Query()
	if code, ok := query["roomCode"]; ok {
		return code[0]
	}
	return ""
}

func (ps *PoolSet) createSocketPoolClient(conn *websocket.Conn, id string, pool remote.ClientPool) *Client {
	client := &Client{
		IDStr:      id,
		Conn:       conn,
		ParentPool: pool,
	}
	pool.AddClient(client)
	return client
}

//BroadcastMsg send a message to all Clients
func (pool *Pool) BroadcastMsg(message remote.Message) {
	for _, client := range pool.Clients {
		if err := client.Send(message); err != nil {
			fmt.Println(err)
		}
	}
}

func (pool *Pool) Count() int {
	return len(pool.Clients)
}

// BroadcastMsgOthers Send to all clients except the sender
func (pool *Pool) BroadcastMsgOthers(message remote.Message, fromID string) {
	for _, client := range pool.Clients {
		if client.ID() != fromID {
			if err := client.Send(message); err != nil {
				fmt.Println(err)
			}
		}
	}
}

// Start runs the event loop
func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.AddClient(client)
			break
		case client := <-pool.Unregister:
			pool.RemoveClient(client)
			break
		case message := <-pool.Messages:
			if message.ToID == "ALL" {
				pool.BroadcastMsgOthers(message, message.FromID)
			} else {
				pool.Clients[message.ToID].Receive(message)
			}
		}
	}
}
