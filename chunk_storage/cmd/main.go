package main

import (
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	proto "str/chunk_storage/grpc"
	serviceimp "str/chunk_storage/internal/service/imp"
	chunkstorage "str/chunk_storage/internal/storage/fs"
	csgrpc "str/chunk_storage/internal/transport/grpc"
)

func main() {
	logger := log.New(os.Stdout, "logger: ", log.Llongfile)

	// Creating chunk storage service
	chunkStorage := chunkstorage.New("/data")
	service := serviceimp.New(chunkStorage)

	// Creating grpc server
	grpcServer := csgrpc.NewServer(service)
	grpcSrv := grpc.NewServer()
	proto.RegisterStorageServiceServer(grpcSrv, grpcServer)

	listener, err := net.Listen("tcp", "0.0.0.0:8000")
	if err != nil {
		logger.Fatal(err)
	}

	logger.Println("running storage...")

	if err := grpcSrv.Serve(listener); err != nil {
		if err = grpcSrv.Serve(listener); err != nil && err != grpc.ErrServerStopped {
			logger.Fatal(err)
		}
	}
}
