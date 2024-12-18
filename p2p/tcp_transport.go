package p2p

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"sync"
)

const NetworkTCP = "tcp"

// TCPPeer represents remote node over a TCP established connection
type TCPPeer struct {
	conn net.Conn
	// connection retrieved by dialing 	 ==> outbound == true
	// connection retrieved by accepting ==> outbound == false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{conn: conn, outbound: outbound}
}

func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

type TCPTransport struct {
	listenAddr string
	listener   net.Listener
	shakeHands HandshakeFunc
	decoder    Decoder
	rpcChanel  chan RPC
	onPeer     func(Peer) error

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransport(listenAddr string, opts ...Opt) *TCPTransport {
	tr := &TCPTransport{
		listenAddr: listenAddr,
		decoder:    &GOBDecoder{},
		shakeHands: NOPHandshakeFunc,
		rpcChanel:  make(chan RPC),
	}

	for _, opt := range opts {
		opt(tr)
	}

	return tr
}

type Opt func(*TCPTransport)

func WithDecoder(d Decoder) Opt {
	return func(transport *TCPTransport) {
		transport.decoder = d
	}
}

func WithHandShaker(hs HandshakeFunc) Opt {
	return func(transport *TCPTransport) {
		transport.shakeHands = hs
	}
}

// Consume implements the transport interface,
// which will return readonly channel received
// from another peer in the network
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcChanel
}

func (t *TCPTransport) ListenAndAccept() error {
	const op = "TCPTransport.ListenAndAccept"

	var err error
	t.listener, err = net.Listen(NetworkTCP, t.listenAddr)
	if err != nil {
		return fmt.Errorf("op: %s, err: %w", op, err)
	}

	go t.startAcceptLoop()
	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	const op = "TCPTransport.startAcceptLoop"
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			slog.Error(fmt.Sprintf("op: %s, "+
				"error accepting new conn , err: %s", op, err))
		}

		slog.Info(fmt.Sprintf("op: %s, "+
			"new incoming connection: %v", op, conn))

		go t.handleAccept(conn)
	}
}

func (t *TCPTransport) handleAccept(conn net.Conn) {
	const op = "TCPTransport.handleAccept"

	var err error

	defer func() {
		slog.Error(fmt.Sprintf("op: %s; dropping peer connection, "+
			"err: %s", op, err))
		_ = conn.Close()
	}()

	peer := NewTCPPeer(conn, true)

	if err = t.shakeHands(peer); err != nil {
		_ = conn.Close()
		slog.Error(fmt.Sprintf("op: %s, err: TCP handshake error", op))
		return
	}

	if t.onPeer != nil && t.onPeer(peer) != nil {
		return
	}

	// Read loop
	msg := &RPC{}
	for {
		if err = t.decoder.Decode(conn, msg); err != nil {
			if errors.Is(err, io.EOF) {
				slog.Info(fmt.Sprintf("op: %s, connection has been closed", op))
				return
			}
			slog.Error(fmt.Sprintf("op: %s, TCP decoding err: %s\n", op, err))
			continue
		}

		msg.From = conn.RemoteAddr()
		t.rpcChanel <- *msg
	}
}
