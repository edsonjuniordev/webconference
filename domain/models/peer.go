package models

import (
	"sync"

	"github.com/gofiber/websocket/v2"
)

type Peer struct {
	Lock       *sync.Mutex
	Connection *websocket.Conn
	ID         string
	Name       string
}

func NewPeer(connection *websocket.Conn, ID, name string) *Peer {
	return &Peer{
		Lock:       &sync.Mutex{},
		Connection: connection,
		ID:         ID,
		Name:       name,
	}
}
