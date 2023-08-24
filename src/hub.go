package main

import (
	"log"
	"sync"
	"time"
)

type hub struct {
	connectionsMx sync.RWMutex
	connections   map[*connection]struct{}
	broadcast     chan Message
	logMx         sync.RWMutex
	log           [][]byte
}

func newHub() *hub {
	h := &hub{
		connectionsMx: sync.RWMutex{},
		broadcast:     make(chan Message),
		connections:   make(map[*connection]struct{}),
	}

	go func() {
		for {
			msg := <-h.broadcast
			h.connectionsMx.RLock()
			for c := range h.connections {
				select {
				case c.send <- msg:
				case <-time.After(1 * time.Second):
					log.Printf("shutting down connection %v", *c)
					h.removeConnection(c)
				}
			}
			h.connectionsMx.RUnlock()
		}
	}()

	return h
}

func (h *hub) addConnection(conn *connection) {
	h.connectionsMx.Lock()
	defer h.connectionsMx.Unlock()
	h.connections[conn] = struct{}{}
}

func (h *hub) removeConnection(conn *connection) {
	h.connectionsMx.Lock()
	defer h.connectionsMx.Unlock()
	if _, ok := h.connections[conn]; ok {
		delete(h.connections, conn)
		close(conn.send)
	}
}
