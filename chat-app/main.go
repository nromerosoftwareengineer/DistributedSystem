
package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type webSocketHandler struct {
	upgrader websocket.Upgrader
}

type Message struct {
	From string `json:"from"`
	To   string `json:"to"`
	Body string `json:"body"`
}

var userId_websocket_map = make(map[string]*websocket.Conn)

func (wsh webSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")
	c, err := wsh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error %s when upgrading connection to websocket", err)
		return
	}
	var msg Message
	userId_websocket_map[userId] = c
	for {
		err := c.ReadJSON(&msg)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				err := c.Close()
				if err != nil {
					return
				}
				log.Println("Connection closed by client")
			} else {
				log.Println("Error reading JSON:", err)
			}
			break
		}
		v := userId_websocket_map[msg.To]
		writeMessage := Message{
			From: userId,
			To:   msg.To,
			Body: msg.Body,
		}
		err = v.WriteJSON(writeMessage)
	}
}

func main() {
	webSocketHandler := webSocketHandler{
		upgrader: websocket.Upgrader{},
	}
	http.Handle("/", webSocketHandler)
	log.Print("Starting server...")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
