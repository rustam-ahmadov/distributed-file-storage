package p2p

import (
	"encoding/gob"
	"fmt"
	"io"
)

type Decoder interface {
	Decode(io.Reader, *RPC) error
}

type GOBDecoder struct{}

func (d GOBDecoder) Decode(r io.Reader, msg *RPC) error {
	const op = "GOBDecoder.Decode"
	if err := gob.NewDecoder(r).Decode(msg); err != nil {
		return fmt.Errorf("op: %s, decoding err: %w", op, err)
	}
	return nil
}

type DefaultDecoder struct {
}

func (d DefaultDecoder) Decode(r io.Reader, msg *RPC) error {
	const op = "Default.Decode"
	buf := make([]byte, 1028)
	n, err := r.Read(buf)
	if err != nil {
		return fmt.Errorf("op: %s, decoding err: %w", op, err)
	}

	msg.Payload = buf[:n]
	return nil
}
