package main

import (
	"bytes"
	"context"
	"crypto/sha256"
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

	// Open the file to send
	file, err := os.Open(fileName)
	if err != nil {
		logger.Fatalf("open file: %v", err)
	}
	defer file.Close()

	// Get sha246 hash
	hasher1 := sha256.New()
	if _, err := io.Copy(hasher1, file); err != nil {
		logger.Fatalf("failed to copy file to hasher: %s", err)
	}

	hash1 := hasher1.Sum(nil)

	// Return the pointer to the beginning of the file
	_, err = file.Seek(0, 0)
	if err != nil {
		logger.Fatalf("seek to the beginning of the file: %v", err)
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

	// Get file stats to get it's size
	fileStat, err := file.Stat()
	if err != nil {
		logger.Fatalf("get file stat: %v", err)
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
	buffer := make([]byte, 1024)
	length1 := 0
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		length1 += 1024
		if length1%(1024*1024) == 0 {
			logger.Printf("Uploaded: %d MB\n", length1/(1024*1024))
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

	// Create resulting file
	resultFile, err := os.Create(fileNameResult)
	if err != nil {
		logger.Fatalf("create resulting file: %v", err)
	}
	defer resultFile.Close()

	// Get first message as metadata
	resd, err := streamd.Recv()
	if err != nil {
		logger.Fatalf("get first message with metadata: %v", err)
	}
	logger.Println(resd.GetSuccess().GetMetadata())

	// Writing content to created resulting file
	length2 := 0
	for {
		resd, err := streamd.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Fatalf("get chunk message: %v", err)
		}

		chunk := resd.GetSuccess().GetChunk().GetContent()
		length2 += len(chunk)
		if length2%(1024*1024) == 0 {
			logger.Printf("Downloaded: %d B, %d MB\n", length2, length2/(1024*1024))
		}

		_, err = resultFile.Write(chunk)
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

	hash2 := hasher2.Sum(nil)

	// Check hashes
	if !bytes.Equal(hash1, hash2) {
		logger.Fatalf("\ntest failed:\noriginal file: %x\nresulting file:%x\n", hash1, hash2)
	}

	logger.Println("test successful")
}
