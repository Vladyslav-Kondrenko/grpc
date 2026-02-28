package main

import (
	"context"
	"log"
	"os"

	products "github.com/Vladyslav-Kondrenko/grpc.git/api/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	addr := os.Getenv("GRPC_SERVER_ADDR")
	if addr == "" {
		addr = "localhost:50051"
	}

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	client := products.NewProductServiceClient(conn)
	ctx := context.Background()

	// Fetch
	resp, err := client.Fetch(ctx, &products.FetchRequest{Url: "http://localhost:8080/products"})
	if err != nil {
		log.Fatalf("Fetch: %v", err)
	}
	log.Printf("Fetch: created=%d updated=%d error=%q", resp.GetCreated(), resp.GetUpdated(), resp.GetError())

	// List
	listResp, err := client.List(ctx, &products.ListRequest{})
	if err != nil {
		log.Fatalf("List: %v", err)
	}
	log.Printf("List: total=%d products", listResp.GetTotal())
	for i, p := range listResp.GetProducts() {
		log.Printf("  [%d] id=%q name=%q price=%d", i+1, p.GetId(), p.GetName(), p.GetPrice())
	}
}
