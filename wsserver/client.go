package wsserver

import (
	"golang.org/x/net/websocket"
	"log"
)

type WSClient struct {
	ws     *websocket.Conn
	doneCh chan bool
	InputChan chan interface{}
	OutputChan chan interface{}
	id 	int
	dead 	bool
}

func NewWSClient(ws *websocket.Conn, id int) *WSClient {

	if ws == nil {
		panic("ws cannot be nil")
	}

	inputCh := make(chan interface{}, 5)
	outputCh := make(chan interface{}, 5)
	doneCh := make(chan bool)

	return &WSClient{ws, doneCh, inputCh, outputCh, id, false}
}

func (client *WSClient) GetId() int {
	return client.id
}

func (client *WSClient) IsAlive() bool {
	return !client.dead
}

// Listen read request via chanel
func (c *WSClient) AwaitDeath() {
	go func() {
		for {
			var msg interface{}
			err := websocket.JSON.Receive(c.ws, &msg)
			if err != nil {
				log.Printf("receive error for %v: %v",c.ws.RemoteAddr(), err)
				c.doneCh <- true
				break
			}
			c.OutputChan <- msg
		}
	}()

	for {
		select {

		// receive done request
		case <-c.doneCh:
			log.Printf("Exiting client %d, bailing", c.id)
			c.dead = true
			close(c.InputChan)
			close(c.OutputChan)
			return

		//process any messages to send to the client
		case msg := <-c.InputChan:
		//send this down the pipe to the client
			err := websocket.JSON.Send(c.ws, msg)
			if(err != nil) {
				log.Printf("Error sending update: %v", err)
			}
		}
	}
}