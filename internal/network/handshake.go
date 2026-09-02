package network

type HelloPayload struct {
	NodeID  string `json:"node_id"`
	Address string `json:"address"`
}
