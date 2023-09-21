package grpc

import (
	"context"
	"fmt"
	"io"
	proto "str/bff/grpc"
	"str/bff/internal/entity"
	"str/bff/internal/service"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type grpcServer struct {
	proto.UnimplementedBFFServiceServer
	service service.Service
}

func NewServer(
	service service.Service,
) proto.BFFServiceServer {
	return &grpcServer{
		service: service,
	}
}

func (s *grpcServer) Upload(stream proto.BFFService_UploadServer) error {
	// Get first message in the stream as metadata
	req, err := stream.Recv()
	if err != nil {
		return err
	}

	metadata := &entity.FileMetadata{
		Name:     req.GetMetadata().GetName(),
		Size:     req.GetMetadata().GetSize(),
		Storages: make(map[entity.FilePart]entity.ConnectionID),
	}

	// Checking for bad data
	if metadata.Name == "" || metadata.Size == 0 {
		stream.SendAndClose(&proto.Upload_Response{
			Response: &proto.Upload_Response_Error{Error: &proto.Error{
				Error:        proto.BFFError_BFF_ERROR_BAD_DATA,
				ErrorMessage: fmt.Sprintf("name: %s, size: %d", metadata.Name, metadata.Size),
			}},
		})
	}

	// This is the channel into which we will stream the file to the service layer.
	fileContent := make(chan []byte)
	// defer close(fileContent)

	var g errgroup.Group
	var uuid uuid.UUID

	// Sending the channel to the service layer
	g.Go(func() error {
		var err error
		uuid, err = s.service.Upload(stream.Context(), fileContent, metadata)
		if err != nil {
			return err
		}

		return nil
	})

	// Streaming uploading file into the channel
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

	// Waiting for the completion of the service layer Upload function
	if err := g.Wait(); err != nil {
		return err
	}

	// Close the stream
	return stream.SendAndClose(&proto.Upload_Response{
		Response: &proto.Upload_Response_Success{
			Success: &proto.Upload_Success{
				Id: uuid.String(),
			},
		},
	})
}

func (s *grpcServer) Download(req *proto.Download_Request, stream proto.BFFService_DownloadServer) error {
	// Parse file id as uuid from request
	fileID, err := uuid.Parse(req.GetId())
	if err != nil {
		return err
	}

	// Download function returns two channels - with file content streaming and with errors
	fileMetadata, fileContent, errorChan, err := s.service.Download(stream.Context(), fileID)
	if err != nil {
		return errors.Wrapf(err, "download file content, file id: %s", fileID)
	}

	// Sending first message to the client with metadata
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
		return errors.Wrapf(err, "send metadata for downloadable file, id: %s, name: %s, size: %d", fileID, fileMetadata.Name, fileMetadata.Size)
	}

	// Streaming file content and listening for errors
	for {
		select {
		case err, ok := <-errorChan:
			if ok {
				return err
			}
		case content, ok := <-fileContent:
			if !ok {
				goto End
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
	}

End:
	return nil
}

// TODO: implement
func (s *grpcServer) Delete(ctx context.Context, req *proto.Delete_Request) (*proto.Delete_Response, error) {
	return nil, nil
}
