package main

import (
	"context"
	"log"
	"net"
	"os"

	products "github.com/Vladyslav-Kondrenko/grpc.git/api/proto"
	"github.com/Vladyslav-Kondrenko/grpc.git/internal/app/server"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func main() {
	_ = godotenv.Load()

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://root:example@localhost:27017/?authSource=admin"
	}

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	coll := client.Database("products").Collection("products")
	svc := server.New(coll)

	grpcServer := grpc.NewServer()

	products.RegisterProductServiceServer(grpcServer, svc)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}
	defer lis.Close()

	log.Println("gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
