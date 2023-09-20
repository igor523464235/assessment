package service

import (
	"context"
	"str/bff/internal/server"
	connectionpool "str/bff/internal/server/connection_pool"

	"github.com/google/uuid"
)

type MetadataStorage interface {
	SaveFileInfo(ctx context.Context, file *server.FileMetadata) error
	GetFileInfo(ctx context.Context, fileID uuid.UUID) (*server.FileMetadata, error)
	// TODO:
	// Delete(ctx context.Context, fileID uuid.UUID) error
}

type ConnectionPool interface {
	IterateConnections() <-chan func() (conn *connectionpool.Connection)
	UpdateFreeSpace(id server.ConnectionID, size uint64)
	GetConnectionByFreeSpace() *connectionpool.Connection
	GetConnection(connectionID server.ConnectionID) *connectionpool.Connection
}
