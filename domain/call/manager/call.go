package manager

import (
	"fmt"

	"github.com/edsonjuniordev/webconference/domain/call/signal"
	"github.com/edsonjuniordev/webconference/domain/models"
	"github.com/edsonjuniordev/webconference/domain/store"
	"github.com/gofiber/websocket/v2"
)

type PeerManager struct {
	Store  store.Store
	Room   *models.Room
	PeerID string
}

func NewPeerManager(store store.Store, room *models.Room) *PeerManager {
	return &PeerManager{Store: store, Room: room}
}

func (p *PeerManager) HandleRequest(req signal.SignalRequest, c *websocket.Conn) error {
	switch req.Type {
	case signal.JoinRequest:
		return p.handlePeerJoin(req, c)
	case signal.SdpOffer, signal.SdpAnswer:
		return p.handleSdp(req)
	case signal.IceCandidate:
		return p.handleIceCandidate(req)
	}

	return nil
}

func (p *PeerManager) handlePeerJoin(req signal.SignalRequest, c *websocket.Conn) error {
	peerID, err := p.Store.RoomCollection.AddPeer(p.Room.ID, req.PeerName, c)
	if err != nil {
		return fmt.Errorf("failed to add peer: %v", err)
	}

	p.PeerID = peerID

	signalResponse := signal.NewSignalResponse(signal.JoinResponse, peerID)

	err = c.WriteJSON(signalResponse)
	if err != nil {
		return fmt.Errorf("failed to response join: %v", err)
	}

	return nil
}

func (p *PeerManager) handleSdp(req signal.SignalRequest) error {
	sdp, ok := req.Payload.(map[string]interface{})
	if !ok {
		return fmt.Errorf("failed to read sdp offer")
	}

	peer, err := p.Store.RoomCollection.GetPeer(p.Room.ID, sdp["destination_id"].(string))
	if err != nil {
		return fmt.Errorf("failed to send sdp to peer: %v", err)
	}

	err = peer.Connection.WriteJSON(req)
	if err != nil {
		return fmt.Errorf("failed to exchange sdp: %v", err)
	}

	return nil
}

func (p *PeerManager) handleIceCandidate(req signal.SignalRequest) error {
	iceCandidate, ok := req.Payload.(map[string]interface{})
	if !ok {
		return fmt.Errorf("failed to read sdp offer")
	}

	peer, err := p.Store.RoomCollection.GetPeer(p.Room.ID, iceCandidate["destination_id"].(string))
	if err != nil {
		return fmt.Errorf("failed to send sdp to peer: %v", err)
	}

	err = peer.Connection.WriteJSON(req)
	if err != nil {
		return fmt.Errorf("failed to exchange sdp: %v", err)
	}

	return nil
}
