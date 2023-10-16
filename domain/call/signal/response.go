package signal

type SignalResponse struct {
	Type    string      `json:"type,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
}

func NewSignalResponse(signalType string, payload interface{}) *SignalResponse {
	return &SignalResponse{Type: signalType, Payload: payload}
}
