package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type connection struct {
	send chan Message
	h    *hub
}

type Message struct {
	FromConnection string
	ChatMessage    string `json:"chat_message"`
}

func (c *connection) reader(wg *sync.WaitGroup, wsConn *websocket.Conn) {
	defer wg.Done()
	for {
		_, msg, err := wsConn.ReadMessage()
		if err != nil {
			break
		}
		var message Message
		err = json.Unmarshal(msg, &message)
		message.FromConnection = wsConn.RemoteAddr().String()
		log.Printf("received message from %v: %v", message.FromConnection, message.ChatMessage)
		c.h.broadcast <- message
	}
}

func (c *connection) writer(wg *sync.WaitGroup, wsConn *websocket.Conn) {
	defer wg.Done()
	for msg := range c.send {
		tmpl := template.Must(template.ParseFiles("resources/views/response.html"))
		buffer := new(bytes.Buffer)
		tmpl.Execute(buffer, msg)
		err := wsConn.WriteMessage(websocket.TextMessage, buffer.Bytes())
		if err != nil {
			break
		}
	}
}

var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

type wsHandler struct {
	h *hub
}

func (wsh wsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	c := &connection{send: make(chan Message, 256), h: wsh.h}
	c.h.addConnection(c)
	defer c.h.removeConnection(c)
	var wg sync.WaitGroup
	wg.Add(2)
	go c.writer(&wg, wsConn)
	go c.reader(&wg, wsConn)
	wg.Wait()
	wsConn.Close()
}
