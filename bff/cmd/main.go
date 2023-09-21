package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	_ "github.com/lib/pq"

	proto "str/bff/grpc"
	connectionpool "str/bff/internal/connection_pool"
	"str/bff/internal/entity"
	serviceimp "str/bff/internal/service/imp"
	"str/bff/internal/storage/postgres"
	bffgrpc "str/bff/internal/transport/grpc"
	csgrpc "str/chunk_storage/grpc"
)

func main() {
	logger := log.New(os.Stdout, "logger: ", log.Llongfile)

	// Getting amount of file storages from env
	storagesStr := os.Getenv("STORAGES")
	storages, err := strconv.Atoi(storagesStr)
	if err != nil {
		logger.Fatalf("invalid amount of storage servers: %v", storagesStr)
	}

	if storages < 1 {
		logger.Fatalf("amount of storages must be greather than 0.")
	}

	logger.Println("storages: ", storages)
	// Creating connection to db
	// TODO: move credentials to envs
	dbconn, err := sql.Open("postgres",
		fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			"db", 5432, "user", "password", "db"),
	)
	if err != nil {
		logger.Fatal(err)
	}
	defer dbconn.Close()

	for {
		if err = dbconn.Ping(); err != nil {
			logger.Println(err, "waiting to connect to to database")
			time.Sleep(time.Second)
			continue
		}

		break
	}

	// Creating db service
	db := postgres.New(dbconn)
	err = db.Init(context.Background())
	if err != nil {
		logger.Fatal(err)
	}

	// Creating pool with grpc connections to chunk storages
	connectionPool := connectionpool.New()

	// TODO: move address and port to envs
	for storageNum := 1; storageNum <= storages; storageNum++ {
		grpconn, err := grpc.Dial(fmt.Sprintf("chunk_storage_%d:8000", storageNum),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
		)
		if err != nil {
			logger.Fatal(err)
		}
		defer grpconn.Close()

		connectionPool.AddConnection(&connectionpool.Connection{
			ID:                   entity.ConnectionID(storageNum),
			FreeSpace:            0,
			StorageServiceClient: csgrpc.NewStorageServiceClient(grpconn),
		})
	}

	// Creating main service with db service and connection pool
	service, err := serviceimp.New(db, connectionPool, logger)
	if err != nil {
		logger.Fatal(err)
	}

	// Creating grpc server with the main service
	grpcServer := bffgrpc.NewServer(service)
	grpcSrv := grpc.NewServer()
	proto.RegisterBFFServiceServer(grpcSrv, grpcServer)

	listener, err := net.Listen("tcp", "0.0.0.0:8000")
	if err != nil {
		logger.Fatal(err)
	}

	logger.Println("running bff...")

	if err := grpcSrv.Serve(listener); err != nil {
		if err = grpcSrv.Serve(listener); err != nil && err != grpc.ErrServerStopped {
			logger.Fatal(err)
		}
	}
}
