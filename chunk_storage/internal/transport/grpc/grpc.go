package grpc

import (
	"context"
	"io"

	proto "str/chunk_storage/grpc"
	"str/chunk_storage/internal/service"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type grpcServer struct {
	proto.UnimplementedStorageServiceServer
	storageService service.Service
}

func NewServer(
	storageService service.Service,
) proto.StorageServiceServer {
	return &grpcServer{
		storageService: storageService,
	}
}

func (s *grpcServer) Save(stream proto.StorageService_SaveServer) error {
	req, err := stream.Recv()
	if err != nil {
		return errors.Wrap(err, "recieve first stream message to save file")
	}

	metadata := req.GetChunkMetadata()

	id, err := uuid.Parse(metadata.GetId())
	if err != nil {
		return errors.Wrap(err, "parse file id from metadata")
	}

	chunkContent := make(chan []byte)
	var g errgroup.Group

	g.Go(func() error {
		err = s.storageService.Save(chunkContent, id, uint8(metadata.GetPart()))
		if err != nil {
			return errors.Wrap(err, "save to storage service")
		}

		return nil
	})

	for {
		req, err = stream.Recv()
		if err == io.EOF {
			close(chunkContent)
			break
		}
		if err != nil {
			return errors.Wrap(err, "recieve chunk message to save file")
		}

		chunkContent <- req.GetChunk().GetContent()
	}

	return stream.SendAndClose(&proto.Save_Response{
		Response: &proto.Save_Response_Success{
			Success: &proto.Save_Success{},
		},
	})
}

func (s *grpcServer) Get(req *proto.Get_Request, stream proto.StorageService_GetServer) error {
	id, err := uuid.ParseBytes([]byte(req.GetId()))
	if err != nil {
		return err
	}

	fileContent, errorChan := s.storageService.Get(id, uint8(req.GetPart()))

	for {
		select {
		case err, ok := <-errorChan:
			if ok {
				return err
			}
		case chunk, ok := <-fileContent:
			if !ok {
				goto End
			}

			if err = stream.Send(&proto.Get_Response{
				Response: &proto.Get_Response_Success{
					Success: &proto.Get_Success{
						Chunk: &proto.Chunk{
							Content: chunk,
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

func (s *grpcServer) Erase(ctx context.Context, req *proto.Erase_Request) (*proto.Erase_Response, error) {
	id, err := uuid.ParseBytes([]byte(req.GetId()))
	if err != nil {
		return &proto.Erase_Response{
			Response: &proto.Erase_Response_Error{
				Error: err.Error(),
			},
		}, nil
	}

	err = s.storageService.Erase(id, uint8(req.GetPart()))
	if err != nil {
		return &proto.Erase_Response{
			Response: &proto.Erase_Response_Error{
				Error: err.Error(),
			},
		}, nil
	}

	return &proto.Erase_Response{
		Response: &proto.Erase_Response_Success{Success: &proto.Erase_Success{}},
	}, nil
}

func (s *grpcServer) GetFreeSpace(ctx context.Context, req *proto.GetFreeSpace_Request) (*proto.GetFreeSpace_Response, error) {
	bytes, err := s.storageService.GetFreeSpace()
	if err != nil {
		return &proto.GetFreeSpace_Response{
			Response: &proto.GetFreeSpace_Response_Error{
				Error: err.Error(),
			},
		}, nil
	}

	return &proto.GetFreeSpace_Response{
		Response: &proto.GetFreeSpace_Response_Success{Success: &proto.GetFreeSpace_Success{
			Bytes: bytes,
		}},
	}, nil
}
