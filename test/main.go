package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	bffgrpc "str/bff/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	logger := log.New(os.Stdout, "logger: ", log.Llongfile)

	fileName := "picture.jpg"
	fileNameResult := "result_file.jpg"

	file, err := os.Open(fileName)
	if err != nil {
		logger.Fatalf("open file: %v", err)
	}
	defer file.Close()

	hasher1 := sha256.New()
	if _, err := io.Copy(hasher1, file); err != nil {
		logger.Fatalf("failed to copy file to hasher: %s", err)
	}

	hash1 := hasher1.Sum(nil)

	fileStat, err := file.Stat()
	if err != nil {
		logger.Fatalf("get file stat: %v", err)
	}

	// Create connection to bff server
	conn, err := grpc.Dial("localhost:8000", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		logger.Fatalf("connect to bff server: %v", err)
	}
	defer conn.Close()

	// Create client with the connection
	client := bffgrpc.NewBFFServiceClient(conn)

	// Get upload stream
	stream, err := client.Upload(context.Background())
	if err != nil {
		logger.Fatalf("upload file: %v", err)
	}

	// Send metadata as the first message
	if err := stream.Send(&bffgrpc.Upload_Request{
		Data: &bffgrpc.Upload_Request_Metadata{
			Metadata: &bffgrpc.FileMetadata{
				Name: fileName,
				Size: uint64(fileStat.Size()),
			},
		},
	}); err != nil {
		logger.Fatalf("send metadata: %v", err)
	}

	// Reading file and sending it's parts
	buffer := make([]byte, 10)
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Fatalf("read file: %v", err)
		}

		if err := stream.Send(&bffgrpc.Upload_Request{
			Data: &bffgrpc.Upload_Request_Chunk{
				Chunk: &bffgrpc.FileChunk{
					Content: buffer[:n],
				},
			},
		}); err != nil {
			logger.Fatalf("send file chunk: %v", err)
		}
	}

	// Close stream
	res, err := stream.CloseAndRecv()
	if err != nil {
		logger.Fatalf("close stream %v", err)
	}

	// Download our uploaded file
	streamd, err := client.Download(context.Background(), &bffgrpc.Download_Request{
		Id: res.GetSuccess().GetId(),
	})
	if err != nil {
		logger.Fatalf("download stream: %v", err)
	}

	resultFile, err := os.Create(fileNameResult)
	if err != nil {
		logger.Fatalf("create resulting file: %v", err)
	}
	defer resultFile.Close()

	resd, err := streamd.Recv()
	if err != nil {
		logger.Fatalf("get first message with metadata: %v", err)
	}
	fmt.Println(resd.GetSuccess().GetMetadata())

	// Writing content to created resulting file
	for {
		resd, err := streamd.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Fatalf("get chunk message: %v", err)
		}

		_, err = resultFile.Write(resd.GetSuccess().GetChunk().GetContent())
		if err != nil {
			logger.Fatalf("write chunk to file: %v", err)
		}
	}

	// Open saved file
	savedFile, err := os.Open(fileNameResult)
	if err != nil {
		logger.Fatalf("failed to open file: %s", err)
	}
	defer savedFile.Close()

	// Get hash of the saved file
	hasher2 := sha256.New()
	if _, err := io.Copy(hasher2, savedFile); err != nil {
		logger.Fatalf("failed to copy file2 to hasher: %s", err)
	}

	hash2 := hasher1.Sum(nil)

	// Check hashes
	if !bytes.Equal(hash1, hash2) {
		logger.Printf("original file: %x\nresulting file:%x\n", hash1, hash2)
	}

	logger.Println("test successful")
}
