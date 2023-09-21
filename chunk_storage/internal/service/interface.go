package service

import "github.com/google/uuid"

type Service interface {
	Save(content <-chan []byte, chunkID uuid.UUID, part uint8) error
	Get(chunkID uuid.UUID, part uint8) (<-chan []byte, <-chan error)
	Erase(chunkID uuid.UUID, part uint8) error
	GetFreeSpace() (uint64, error)
}
