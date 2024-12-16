package client

import (
	"log"
	
	pb "api-service/pkg/proto"

	"google.golang.org/grpc"
)

func NewDBServiceClient() (pb.DBServiceClient, *grpc.ClientConn) {
	dbConn, err := grpc.Dial("db-service:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to db-service: %v", err)
	}

	return pb.NewDBServiceClient(dbConn), dbConn
}
