package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Msg struct {
	clientKey string
	text      string
}

type NewClientEvent struct {
	clientKey string
	msgChain  chan Msg
}

var (
	clientRequest    = make(chan *NewClientEvent, 100)
	clientDisconnect = make(chan string, 100)
	broadcast        = make(chan *Msg, 100)
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// // Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// serveWs handles websocket requests from the peer.
func handlesWebSocketRequests(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	conn.WriteMessage(websocket.BinaryMessage, []byte{'H', 'e', 'l', 'l', 'o', '\n'})

	msgChan := make(chan Msg, 100)
	clientKey := r.RemoteAddr
	clientRequest <- &NewClientEvent{clientKey: clientKey, msgChain: msgChan}

	defer func(clientKey string) {
		clientDisconnect <- clientKey
		conn.Close()
	}(clientKey)

	go func() {
		for msg := range msgChan {
			conn.WriteMessage(websocket.BinaryMessage, []byte(msg.text))
		}
	}()

	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	go func() {
		ticker := time.NewTicker(pingPeriod)
		defer func() {
			ticker.Stop()
		}()

		for {
			select {
			case <-ticker.C:
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			}
		}
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		broadcast <- &Msg{text: string(message), clientKey: clientKey} // ?
	}
}

func router() {
	clients := make(map[string]chan Msg)
	for {
		select {
		case req := <-clientRequest:
			clients[req.clientKey] = req.msgChain
			log.Println("WebSocket connected: " + req.clientKey)
		case clientKey := <-clientDisconnect:
			close(clients[clientKey])
			delete(clients, clientKey)
			log.Println("WebSocket disconnected: " + clientKey)
		case msg := <-broadcast:
			for _, msgChan := range clients {
				if cap(msgChan) > len(msgChan) {
					msgChan <- *msg
				}
			}
		}

	}
}
