![Go Report](https://goreportcard.com/badge/github.com/Vladyslav-Kondrenko/grpc)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/Vladyslav-Kondrenko/grpc)
![GitHub repo size](https://img.shields.io/github/repo-size/Vladyslav-Kondrenko/grpc)
![GitHub last commit](https://img.shields.io/github/last-commit/Vladyslav-Kondrenko/grpc)
![License](https://img.shields.io/badge/license-MIT-green)

<img align="right" width="220" src="image.png">

# Go gRPC Products Service

gRPC service for managing products: fetch from CSV, store in MongoDB, list with pagination and sorting.

## Overview

- **Fetch** — download products from a CSV URL and upsert into MongoDB.
- **List** — read products with pagination and sorting (by `name` or `price`).

Implemented in Go with protocol buffers. Server listens on `:50051`; client calls Fetch then List.

## Solution notes

- :zap: gRPC API with protobuf (`api/proto/products.proto`)
- :floppy_disk: MongoDB for product storage; optional Postgres (see `deployments/`)
- :whale: Docker Compose for MongoDB (and other services)
- :file_folder: Standard layout: `cmd/server`, `cmd/client`, `internal/app/server`

## HOWTO

```bash
# Start dependencies (MongoDB, etc.)
docker compose up -d

# Run gRPC server (default :50051)
go run ./cmd/server

# In another terminal — run client
go run ./cmd/client
```

Env: `MONGODB_URI` (server), `GRPC_SERVER_ADDR` (client, default `localhost:50051`).

## gRPC API

```proto
service ProductService {
  rpc Fetch(FetchRequest) returns (FetchResponse);
  rpc List(ListRequest) returns (ListResponse);
}
```

- **Fetch**: `url` → downloads `;`-separated CSV, upserts by product name; returns `created` / `updated` counts.
- **List**: `paging` + `sorting` → returns `products` and `total`.

## Project structure

```
├── api/proto/          # .proto + generated Go
├── cmd/server/         # gRPC server
├── cmd/client/         # gRPC client
├── internal/app/server # ProductService implementation
├── deployments/        # DB init (e.g. Postgres)
├── docker-compose.yml
└── go.mod
```

## Regenerating protobuf

```bash
protoc --go_out=. --go-grpc_out=. api/proto/products.proto
```
