package signal

type SignalRequest struct {
	Type     string      `json:"type,omitempty"`
	PeerName string      `json:"peer_name,omitempty"`
	Payload  interface{} `json:"payload,omitempty"`
}
