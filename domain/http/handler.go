package http

import (
	"log"
	"net/http"

	"github.com/edsonjuniordev/webconference/domain/call/manager"
	"github.com/edsonjuniordev/webconference/domain/call/signal"
	"github.com/edsonjuniordev/webconference/domain/store"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type Handler struct {
	Store store.Store
}

func (h Handler) Register(app *fiber.App) {
	app.Post("/room", h.create)
	app.Get("ws/room/:room_id", websocket.New(h.join))
}

func (h Handler) create(ctx *fiber.Ctx) error {
	roomID := h.Store.RoomCollection.Create()

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"room_id": roomID,
	})
}

func (h Handler) join(c *websocket.Conn) {
	roomID := c.Params("room_id")

	room, err := h.Store.RoomCollection.Get(roomID)
	if err != nil {
		log.Printf("failed to join the room: %v", err)
		notFound := signal.NewSignalResponse(signal.RoomNotFound, "")
		c.Conn.WriteJSON(notFound)
		c.Conn.Close()
		return
	}

	peerManager := manager.NewPeerManager(h.Store, room)

	for {
		var req signal.SignalRequest

		err := c.ReadJSON(&req)
		if err != nil {
			log.Printf("failed to read json %v", err)
			room.OnLeave(peerManager.PeerID)
			err := h.Store.RoomCollection.DelPeer(roomID, peerManager.PeerID)
			log.Printf("failed to remove peer: %v", err)
			break
		}
	}
}
