package imp

import (
	"context"

	"github.com/google/uuid"

	connectionpool "str/bff/internal/connection_pool"
	"str/bff/internal/entity"
)

type MetadataStorage interface {
	SaveFileInfo(ctx context.Context, file *entity.FileMetadata) error
	GetFileInfo(ctx context.Context, fileID uuid.UUID) (*entity.FileMetadata, error)
	// TODO:
	// Delete(ctx context.Context, fileID uuid.UUID) error
}

type ConnectionPool interface {
	IterateConnections() <-chan func() (conn *connectionpool.Connection)
	UpdateFreeSpace(id entity.ConnectionID, size uint64)
	GetConnectionByFreeSpace() *connectionpool.Connection
	GetConnection(connectionID entity.ConnectionID) *connectionpool.Connection
}
