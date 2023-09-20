package server

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	Upload(ctx context.Context, c <-chan []byte, fileMetadata *FileMetadata) (uuid.UUID, error)
	Download(ctx context.Context, fileID uuid.UUID) (*FileMetadata, <-chan []byte, error)
	// TODO:
	// Delete(chunkID uuid.UUID, part uint8) error
	// GetFreeSpace() (uint64, error)
}
