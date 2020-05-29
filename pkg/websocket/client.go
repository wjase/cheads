package websocket

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/wjase/cheads/pkg/remote"
)

// Client A client connection with an ID
type Client struct {
	IDStr      string
	Conn       *websocket.Conn
	ParentPool remote.ClientPool
	Listener   remote.MessageListener
	// Messages chan Message
}

// ReadLoop runs forever passing on messages to the pool
func (c *Client) ReadLoop() {
	defer func() {
		c.Close()
	}()

	for {
		_, msgToSendBytes, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		msg := remote.Message{}
		json.NewDecoder(strings.NewReader(string(msgToSendBytes))).Decode(&msg)
		c.ParentPool.Dispatch(msg)
	}
}

// ID the unique id of the client
func (c *Client) ID() string {
	return c.IDStr
}

// Close closes the connection and disconnects from the pool
func (c *Client) Close() {
	c.ParentPool.RemoveClient(c)
	c.Conn.Close()
}

// Subscribe subscribes to inbound messages
func (c *Client) Subscribe(listener remote.MessageListener) {
	c.Listener = listener
}

//Send Sends a message to the client
func (c *Client) Send(message remote.Message) error {
	return c.Conn.WriteJSON(message)
}

// Receive process a message on this clients listener
func (c *Client) Receive(message remote.Message) error {
	return c.Listener.OnMessage(message)
}

func (c *Client) Pool() remote.ClientPool {
	return c.ParentPool
}
