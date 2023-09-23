package service

import (
	"context"
	"str/bff/internal/entity"

	"github.com/google/uuid"
)

type Service interface {
	Upload(ctx context.Context, c <-chan []byte, fileMetadata *entity.FileMetadata) (uuid.UUID, error)
	Download(ctx context.Context, fileID uuid.UUID) (*entity.FileMetadata, <-chan []byte, <-chan error, error)
	UpdateFreeSpaces(ctx context.Context) error
	// TODO:
	// Delete(chunkID uuid.UUID, part uint8) error
}
