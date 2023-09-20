package chunkstorage

import (
	"github.com/google/uuid"
)

type Service interface {
	Save(chunkContent <-chan []byte, chunkID uuid.UUID, part uint8) error
	Get(chunkID uuid.UUID, part uint8) ([]byte, error)
	GetFreeSpace() (uint64, error)
	Erase(chunkID uuid.UUID, part uint8) error
}
