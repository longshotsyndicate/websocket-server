package main

import (
	"github.com/longshotsyndicate/websocket-server/wsserver"
	"log"
	"golang.org/x/net/websocket"
	"net/http"
	"time"
)

func main() {
	log.Printf("ALIVE")


	// websocket handler
	chans := make(map[int]chan interface{})
	clientIds := 0;
	onConnected := func(ws *websocket.Conn) {
		defer ws.Close()
		clientId := clientIds
		clientIds++
		client := wsserver.NewWSClient(ws, clientId)
		chans[clientId] = client.InputChan
		client.AwaitDeath()
		delete(chans, clientId)

	}

	http.Handle("/ws", websocket.Handler(onConnected))


	//send a message to everyone
	go func() {
		for {
			time.Sleep(1*time.Second)
			log.Println("hello")
			for _, clientChan := range chans {
				clientChan <- dave{"howdy", "hi", 45645}
			}
		}
	}()



	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}

}

type dave struct {
	Dave string
	Mark string
	Alessio int
}



