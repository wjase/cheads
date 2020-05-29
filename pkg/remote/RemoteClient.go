package remote

//go:generate mockgen -destination=../mock_remote/MockClient.go	-package=mock_remote . Client

// MessageListener responds to asyn mesages
type MessageListener interface {
	//OnMessage Called when a message is recieved
	OnMessage(message Message) error
}

// Client A generic interface for a connected messaging client
type Client interface {
	// ID the unique id of the client
	ID() string

	// Subscribe subscribes to inbound messages
	Subscribe(listener MessageListener)

	//Send Sends a message to the client
	Send(message Message) error

	//Receive Receive a message on the client
	Receive(message Message) error

	// Pool The pool managing client
	Pool() ClientPool

	//Close closes the connection
	Close()
}
