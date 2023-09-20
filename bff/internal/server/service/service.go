package service

import (
	"context"
	"io"
	"log"
	"str/bff/internal/server"
	connectionpool "str/bff/internal/server/connection_pool"
	csgrpc "str/storage/grpc"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Service struct {
	storage        MetadataStorage
	connection     *connectionpool.Connection
	connectionPool ConnectionPool
	logger         *log.Logger
}

func New(storage MetadataStorage, connectionPool *connectionpool.ConnectionPool, logger *log.Logger) (*Service, error) {
	s := &Service{
		connectionPool: connectionPool,
		storage:        storage,
		logger:         logger,
	}

	// During creation, the service updates data about available space on all known file storages to it.
	err := s.UpdateFreeSpaces(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "update free spaces of storages for service")
	}

	return s, nil
}

func (s *Service) Upload(ctx context.Context, fileContent <-chan []byte, fileMetadata *server.FileMetadata) (uuid.UUID, error) {
	partSize := fileMetadata.Size / 6
	var contentSize uint64
	part := server.FilePart(1)
	newConnection := true
	fileMetadata.ID = uuid.New()

	var stream csgrpc.StorageService_SaveClient
	var err error

	for {
		content, ok := <-fileContent
		if !ok {
			break
		}

		if newConnection {
			s.connection = s.connectionPool.GetConnectionByFreeSpace()
			fileMetadata.Storages[part] = s.connection.ID

			stream, err = s.connection.Save(ctx)
			if err != nil {
				return uuid.Nil, err
			}

			// Sending metadata as first message
			err = stream.Send(&csgrpc.Save_Request{
				Data: &csgrpc.Save_Request_ChunkMetadata{
					ChunkMetadata: &csgrpc.ChunkMetadata{
						Id:   fileMetadata.ID.String(),
						Part: uint32(part),
					},
				},
			})
			if err != nil {
				return uuid.Nil, err
			}
			newConnection = false
		}

		err = stream.Send(&csgrpc.Save_Request{
			Data: &csgrpc.Save_Request_Chunk{
				Chunk: &csgrpc.Chunk{
					Content: content,
				},
			},
		})
		if err != nil {
			return uuid.Nil, errors.Wrap(err, "send upload chunk")
		}
		contentSize += uint64(len(content))
		if contentSize >= partSize {
			s.connectionPool.UpdateFreeSpace(s.connection.ID, s.connection.FreeSpace-contentSize)
			contentSize = 0
			part++
			newConnection = true
			_, err = stream.CloseAndRecv()
			if err != nil {
				return uuid.Nil, errors.Wrap(err, "close upload stream")
			}
		}
	}

	if !newConnection {
		_, err = stream.CloseAndRecv()
		if err != nil {
			return uuid.Nil, errors.Wrap(err, "close final upload stream")
		}
	}

	// Save uploaded file info to the database
	err = s.storage.SaveFileInfo(ctx, fileMetadata)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "save file info")
	}

	return fileMetadata.ID, nil
}

func (s *Service) Download(ctx context.Context, fileID uuid.UUID) (*server.FileMetadata, <-chan []byte, error) {
	fileContent := make(chan []byte)

	fileMetadata, err := s.storage.GetFileInfo(ctx, fileID)
	if err != nil {
		return nil, fileContent, errors.Wrap(err, "get file info to download file")
	}

	go func() {
		defer close(fileContent)

		for _, storage := range fileMetadata.Storages.SortByParts() {
			connection := s.connectionPool.GetConnection(storage.ConnectionID)

			stream, err := connection.Get(ctx, &csgrpc.Get_Request{
				Id:   fileID.String(),
				Part: uint32(storage.Part),
			})
			if err != nil {
				s.logger.Println("connection.Connection.Get", err)
				// return nil, fileContent, err
				return
			}

			for {
				resp, err := stream.Recv()
				if err == io.EOF {
					// Закрыть поток после завершения передачи
					if err := stream.CloseSend(); err != nil {
						s.logger.Println("close stream", err)
						// return nil, fileContent, err
						return
					}
					break
				}
				if err != nil {
					s.logger.Println("Get stream.Recv()", err)
					return
					// return nil, fileContent, err
				}

				fileContent <- resp.GetSuccess().GetChunk().GetContent()
			}
		}
	}()

	return fileMetadata, fileContent, nil
}

// TODO: implement
func (s *Service) Delete(id uuid.UUID) (uuid.UUID, error) {
	return uuid.Nil, nil
}

// UpdateFreeSpaces updates information about available space on all file storages in the connection pool.
func (s *Service) UpdateFreeSpaces(ctx context.Context) error {
	for ch := range s.connectionPool.IterateConnections() {
		connection := ch()

		res, err := connection.GetFreeSpace(ctx, &csgrpc.GetFreeSpace_Request{})
		if err != nil {
			return errors.Wrap(err, "get storage free space")
		}

		s.connectionPool.UpdateFreeSpace(connection.ID, res.GetSuccess().Bytes)
	}

	return nil
}
