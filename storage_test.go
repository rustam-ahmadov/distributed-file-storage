package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	r := bytes.NewReader([]byte("hello"))
	path, hash, _ := CASPathTransformFunc(r)
	assert.Equal(t, len(path), 9+2)
	assert.Equal(t, len(hash), 64)
}

func TestStorage(t *testing.T) {
	s := NewStorage(WithPathTransformer(CASPathTransformFunc))

	r := bytes.NewReader([]byte("some jpg picture"))
	if err := s.writeStream(r); err != nil {
		t.Error(err)
	}
}
