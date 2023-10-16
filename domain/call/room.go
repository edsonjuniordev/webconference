package call

import (
	"log"

	"github.com/edsonjuniordev/webconference/domain/call/signal"
	"github.com/edsonjuniordev/webconference/domain/models"
	"github.com/gofiber/fiber/v2"
)

func RoomOnJoin(room *models.Room) func(peerID string) {
	return func(peerID string) {
		peerName := room.Peers[peerID].Name
		log.Printf("onJoin: %s %s", peerID, peerName)

		resp := signal.NewSignalResponse(signal.NewPeer, fiber.Map{
			"name": peerName,
			"id":   peerID,
		})

		broadcastSignal(room, peerID, resp)
	}
}

func RoomOnLeave(room *models.Room) func(peerID string) {
	return func(peerID string) {
		resp := signal.NewSignalResponse(signal.PeerLeave, peerID)
		broadcastSignal(room, peerID, resp)
	}
}

func RoomOnClose(room *models.Room) func() {
	return func() {
		resp := signal.NewSignalResponse(signal.RoomClose, nil)
		broadcastSignal(room, "", resp)
	}
}

func broadcastSignal(room *models.Room, exceptID string, signal *signal.SignalResponse) {
	for id, peer := range room.Peers {
		if id == exceptID {
			continue
		}

		peer.Lock.Lock()
		if err := peer.Connection.WriteJSON(signal); err != nil {
			log.Printf("failed to send signal of type %s to %s: %v", signal.Type, id, err)
		}
		peer.Lock.Unlock()
	}
}
