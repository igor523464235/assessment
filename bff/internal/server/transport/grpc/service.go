package service

import (
	"context"
	"fmt"
	"io"
	proto "str/bff/grpc"
	server "str/bff/internal/server"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type grpcServer struct {
	proto.UnimplementedBFFServiceServer
	service server.Service
}

func NewServer(
	service server.Service,
) proto.BFFServiceServer {
	return &grpcServer{
		service: service,
	}
}

func (s *grpcServer) Upload(stream proto.BFFService_UploadServer) error {
	req, err := stream.Recv()
	if err != nil {
		return err
	}

	metadata := &server.FileMetadata{
		Name:     req.GetMetadata().GetName(),
		Size:     req.GetMetadata().GetSize(),
		Storages: make(map[server.FilePart]server.ConnectionID),
	}

	if metadata.Name == "" || metadata.Size == 0 {
		stream.SendAndClose(&proto.Upload_Response{
			Response: &proto.Upload_Response_Error{Error: &proto.Error{
				Error:        proto.BFFError_BFF_ERROR_BAD_DATA,
				ErrorMessage: fmt.Sprintf("name: %s, size: %d", metadata.Name, metadata.Size),
			}},
		})
	}

	fileContent := make(chan []byte)
	var g errgroup.Group

	var uuid uuid.UUID
	g.Go(func() error {
		var err error
		uuid, err = s.service.Upload(stream.Context(), fileContent, metadata)
		if err != nil {
			return err
		}

		return nil
	})

	for {
		req, err = stream.Recv()
		if err == io.EOF {
			close(fileContent)
			break
		}
		if err != nil {
			return err
		}

		chunk := req.GetChunk()
		content := chunk.GetContent()

		fileContent <- content
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return stream.SendAndClose(&proto.Upload_Response{
		Response: &proto.Upload_Response_Success{
			Success: &proto.Upload_Success{
				Id: uuid.String(),
			},
		},
	})
}

func (s *grpcServer) Download(req *proto.Download_Request, stream proto.BFFService_DownloadServer) error {
	fileID, err := uuid.Parse(req.GetId())
	if err != nil {
		return err
	}

	fileMetadata, fileContent, err := s.service.Download(stream.Context(), fileID)
	if err != nil {
		return errors.Wrapf(err, "download file content, file id: %s", fileID)
	}

	if err = stream.Send(&proto.Download_Response{
		Response: &proto.Download_Response_Success{
			Success: &proto.Download_Success{
				Data: &proto.Download_Success_Metadata{
					Metadata: &proto.FileMetadata{
						Name: fileMetadata.Name,
						Size: fileMetadata.Size,
					},
				},
			},
		},
	}); err != nil {
		return errors.Wrapf(err, "dend metadata for downloadable file, id: %s, name: %s, size: %d", fileID, fileMetadata.Name, fileMetadata.Size)
	}

	for {
		content, ok := <-fileContent
		if !ok {
			break
		}

		if err = stream.Send(&proto.Download_Response{
			Response: &proto.Download_Response_Success{
				Success: &proto.Download_Success{
					Data: &proto.Download_Success_Chunk{
						Chunk: &proto.FileChunk{
							Content: content,
						},
					},
				},
			},
		}); err != nil {
			return err
		}
	}

	return nil
}

func (s *grpcServer) Delete(ctx context.Context, req *proto.Delete_Request) (*proto.Delete_Response, error) {
	return nil, nil
}
