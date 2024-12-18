package p2p

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCPTransport(t *testing.T) {
	listenAddr := ":8080"
	tr := NewTCPTransport(listenAddr)

	assert.Equal(t, tr.listenAddr, listenAddr)
	assert.Nil(t, tr.ListenAndAccept())
}
