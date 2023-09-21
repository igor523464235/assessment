package imp

import (
	"github.com/google/uuid"
)

type Storage interface {
	Save(chunkContent <-chan []byte, hunkID uuid.UUID, part uint8) error
	Get(chunkID uuid.UUID, part uint8) (<-chan []byte, <-chan error)
	Erase(chunkID uuid.UUID, part uint8) error
}
