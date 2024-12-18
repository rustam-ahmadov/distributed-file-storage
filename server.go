package main

import (
	"github.com/rustam-ahmadov/distributed-file-storage/p2p"
)

type Storage interface {
}

type FileServer struct {
	storage   Storage
	transport p2p.Transport
}

func NewFileServer(storage Storage, transport p2p.Transport) *FileServer {
	return &FileServer{
		storage:   storage,
		transport: transport,
	}
}

func (f *FileServer) Start() {

}
