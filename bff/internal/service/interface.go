package service

import (
	"context"
	"str/bff/internal/entity"

	"github.com/google/uuid"
)

type Service interface {
	Upload(ctx context.Context, c <-chan []byte, fileMetadata *entity.FileMetadata) (uuid.UUID, error)
	Download(ctx context.Context, fileID uuid.UUID) (*entity.FileMetadata, <-chan []byte, <-chan error, error)
	// TODO:
	// Delete(chunkID uuid.UUID, part uint8) error
	// GetFreeSpace() (uint64, error)
}
