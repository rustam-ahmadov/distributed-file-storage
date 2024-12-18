package p2p

// Peer is an interface that represent remote node
type Peer interface {
	Close() error
}

// Transport is anything that handles communication
// between the nodes in the network, this can be of the
// form (TCP, UDP, WEBSOCKETS ...)
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
}
