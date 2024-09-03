package network

// Message is received from peers in a p2p network.
type Message struct {
	ID     string
	PeerID string
	Data   []byte
}
