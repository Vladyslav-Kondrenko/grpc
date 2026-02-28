package server

import (
	"context"
	"encoding/csv"
	"io"
	"net/http"
	"strconv"

	products "github.com/Vladyslav-Kondrenko/grpc.git/api/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Server реализует gRPC-сервис продуктов (Fetch, List).
type Server struct {
	products.UnimplementedProductServiceServer
	coll *mongo.Collection
}

// New создаёт сервер с доступом к коллекции products в MongoDB.
func New(coll *mongo.Collection) *Server {
	return &Server{coll: coll}
}

// Fetch скачивает CSV по URL, парсит и создаёт/обновляет продукты в Mongo.
func (s *Server) Fetch(ctx context.Context, req *products.FetchRequest) (*products.FetchResponse, error) {
	url := req.GetUrl()
	if url == "" {
		return &products.FetchResponse{Error: "url is required"}, nil
	}

	httpResp, err := http.Get(url)
	if err != nil {
		return &products.FetchResponse{Error: err.Error()}, nil
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return &products.FetchResponse{Error: "HTTP " + strconv.Itoa(httpResp.StatusCode)}, nil
	}

	reader := csv.NewReader(httpResp.Body)
	reader.Comma = ';'

	_, err = reader.Read()
	if err != nil {
		return &products.FetchResponse{Error: "read header: " + err.Error()}, nil
	}

	var created, updated int32
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return &products.FetchResponse{Error: "read row: " + err.Error(), Created: created, Updated: updated}, nil
		}
		if len(row) < 2 {
			continue
		}
		name := row[0]
		priceStr := row[1]
		price, err := strconv.Atoi(priceStr)
		if err != nil {
			continue
		}

		doc := bson.M{"_id": name, "name": name, "price": price}
		result, err := s.coll.ReplaceOne(ctx, bson.M{"_id": name}, doc, options.Replace().SetUpsert(true))
		if err != nil {
			return &products.FetchResponse{Error: err.Error(), Created: created, Updated: updated}, nil
		}
		if result.UpsertedCount > 0 {
			created++
		}
		if result.ModifiedCount > 0 {
			updated++
		}
	}

	return &products.FetchResponse{Created: created, Updated: updated}, nil
}

// List возвращает продукты из Mongo с пагинацией и сортировкой.
func (s *Server) List(ctx context.Context, req *products.ListRequest) (*products.ListResponse, error) {
	page := int32(1)
	pageSize := int32(10)
	if p := req.GetPaging(); p != nil {
		page = p.GetPage()
		pageSize = p.GetPageSize()
	}
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	field := "name"
	asc := 1
	if sort := req.GetSorting(); sort != nil {
		field = sort.GetField()
		if !sort.GetAscending() {
			asc = -1
		}
	}
	if field != "name" && field != "price" {
		field = "name"
	}

	total, err := s.coll.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	opts := options.Find().
		SetSkip(int64((page - 1) * pageSize)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: field, Value: asc}})

	cursor, err := s.coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var list []*products.Product
	for cursor.Next(ctx) {
		var doc struct {
			ID    string `bson:"_id"`
			Name  string `bson:"name"`
			Price int32  `bson:"price"`
		}
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		list = append(list, &products.Product{Id: doc.ID, Name: doc.Name, Price: doc.Price})
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &products.ListResponse{Products: list, Total: int32(total)}, nil
}
