package controller

import (
	"fmt"
	"log"
	"model"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan *model.Offer)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func RootHandler(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "index.html")
}

func Writer(offer *model.Offer) {
	broadcast <- offer
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// register client
	clients[ws] = true
}

func Echo() {
	for {
		val := <-broadcast
		latlong := fmt.Sprintf("%f %f %s", val.Bid_Price, val.Bid_Price, val.Title)
		// send to every client that is currently connected
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(latlong))
			if err != nil {
				log.Printf("Websocket error: %s", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
