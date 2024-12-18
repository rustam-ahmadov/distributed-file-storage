package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

type PathTransformer func(r io.Reader) (storagePath string, fileName string, err error)

const (
	bufLength              = 1024
	dirLength              = 3
	hexEncodedSha256Length = 64
	dirLevel               = 3
)

// CASPathTransformFunc is func that calculates path
// based on the content itself, providing a digital
// fingerprint that ensures the data's authenticity and uniqueness.
func CASPathTransformFunc(r io.Reader) (storagePath string, contentHash string, err error) {
	const op = "storage.CASPathTransformFunc"
	hash := sha256.New()
	_, err = io.Copy(hash, r)
	if err != nil {
		return "", "", fmt.Errorf("op: %s, err: %w", op, err)
	}

	hashResult := hash.Sum(nil)
	// []byte{0, 1, 2} -> 0 = 48(ascii) -> 30(hex-decimal) cause 16 * 3 = 48
	hexEncodedHash := hex.EncodeToString(hashResult)

	var path string

	f, s := 0, dirLength
	for range dirLevel {
		path = filepath.Join(path, hexEncodedHash[f:s])
		f, s = s, s+dirLength
	}

	return path, hexEncodedHash, nil
}

type CasStorage struct {
	pathTransformer PathTransformer
}

func NewStorage(opts ...Opt) *CasStorage {
	s := &CasStorage{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

type Opt func(*CasStorage)

func WithPathTransformer(pt PathTransformer) Opt {
	return func(storage *CasStorage) {
		storage.pathTransformer = pt
	}
}

func (s *CasStorage) writeStream(r io.Reader) error {
	const op = "storage.writeStream"
	r.Read()
	path, fileName, err := s.pathTransformer(r)
	if err != nil {
		return fmt.Errorf("op: %s, err: %w", op, err)
	}

	if err = os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}

	finalPath := filepath.Join(path, fileName)

	file, err := os.Create(finalPath)
	defer func() { _ = file.Close() }()
	if err != nil {
		return fmt.Errorf("op: %s, err: %w", op, err)
	}

	n, err := io.Copy(file, r)
	if err != nil {
		return fmt.Errorf("op: %s, err: %w", op, err)
	}

	slog.Info(fmt.Sprintf("op: %s; written bytes: %d", op, n))

	return nil
}
