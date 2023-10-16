package collection

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/edsonjuniordev/webconference/domain/call"
	"github.com/edsonjuniordev/webconference/domain/models"
	"github.com/gofiber/websocket/v2"
)

const (
	RoomIDLength = 8
	PeerIDLength = 16
)

type Room struct {
	Rooms map[string]*models.Room
}

func (r *Room) Create() string {
	roomID := generateRandomHash(RoomIDLength)
	room := models.NewRoom(roomID)
	r.Rooms[roomID] = room

	room.OnJoin = call.RoomOnJoin(room)
	room.OnLeave = call.RoomOnLeave(room)
	room.OnClose = call.RoomOnClose(room)

	return roomID
}

func (r *Room) Get(roomID string) (*models.Room, error) {
	room := r.Rooms[roomID]

	if room == nil {
		return nil, errors.New("room does not exists")
	}

	return room, nil
}

func (r *Room) Del(roomID string) error {
	room, err := r.Get(roomID)
	if err != nil {
		return fmt.Errorf("failed to find the room: %w", err)
	}

	for _, peer := range room.Peers {
		peer.Connection.Close()
	}

	room.OnClose()
	delete(r.Rooms, roomID)

	return nil
}

func (r *Room) AddPeer(roomID, peerName string, conn *websocket.Conn) (string, error) {
	room, err := r.Get(roomID)
	if err != nil {
		return "", fmt.Errorf("failed to find the room: %w", err)
	}

	peerID := generateRandomHash(PeerIDLength)
	peer := models.NewPeer(conn, peerID, peerName)
	room.Peers[peerID] = peer
	room.OnJoin(peerID)

	return peerID, nil
}

func (r *Room) GetPeer(roomID, peerID string) (*models.Peer, error) {
	room, err := r.Get(roomID)
	if err != nil {
		return nil, fmt.Errorf("failed to find the room: %w", err)
	}

	peer := room.Peers[peerID]
	if peer == nil {
		return nil, fmt.Errorf("peer does not exists")
	}

	return peer, nil
}

func generateRandomHash(length int) string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]rune, length)

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
