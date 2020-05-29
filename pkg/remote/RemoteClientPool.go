package remote

// PoolAddedFn callback when a pool is added to a set
type PoolAddedFn func(ClientPool)

// ClientJoinedFn event callback when Client Joined
type ClientJoinedFn func(Client)

// ClientLeftFn event callback when Client Left
type ClientLeftFn func(Client)

// ClientPool generic client pool
type ClientPool interface {
	ID() string
	AddClient(Client)
	RemoveClient(Client)
	Count() int
	SubscribeJoined(ClientJoinedFn)
	Dispatch(Message)
}
