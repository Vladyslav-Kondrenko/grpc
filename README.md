<p align="center">
  <b>Go gRPC Products Service</b> 🐹
</p>

---

## Overview

<table>
  <tr>
    <td valign="top">
      <p>
        This project is a gRPC-based service for managing products backed by MongoDB.
        It provides two main operations:
      </p>
      <ul>
        <li><b>Fetch</b>: download products from a CSV endpoint and upsert them into MongoDB.</li>
        <li><b>List</b>: read products from MongoDB with pagination and sorting.</li>
      </ul>
      <p>
        The service is implemented in Go and uses protocol buffers for API contracts.
      </p>
    </td>
    <td align="right" valign="top" width="220">
      <img src="image.png" alt="Go gRPC Products Service" width="200" />
    </td>
  </tr>
</table>

## Architecture

- **gRPC server** (`cmd/server` + `internal/app/server`):
  - Exposes the `ProductService` with `Fetch` and `List` RPCs.
  - Connects to MongoDB and stores products in the `products.products` collection.
- **gRPC client** (`cmd/client`):
  - Calls `Fetch` with a configurable CSV URL.
  - Calls `List` and prints products to stdout.
- **Protocol buffers** (`api/proto/products.proto`):
  - Defines messages and RPCs for `ProductService`.

## Prerequisites

- **Go** ≥ 1.24
- **Docker & Docker Compose** (recommended for MongoDB)
- **protoc** and Go plugins (only if you need to regenerate `.pb.go` files)

## Running with Docker Compose

From the project root:

```bash
docker compose up -d
```

This will start MongoDB (and any other defined services).  
By default the gRPC server expects MongoDB at:

```text
mongodb://root:example@localhost:27017/?authSource=admin
```

You can override this with the `MONGODB_URI` environment variable.

## Running the gRPC Server

```bash
go run ./cmd/server
```

Environment variables:

- **`MONGODB_URI`**: MongoDB connection string (optional, has a sensible default).

The server listens on **`:50051`**.

## Running the gRPC Client

In a separate terminal, with the server running:

```bash
go run ./cmd/client
```

Environment variables:

- **`GRPC_SERVER_ADDR`**: gRPC server address (default: `localhost:50051`).

The client:

1. Calls **`Fetch`** with the URL `http://localhost:8080/products` (can be changed in code).
2. Calls **`List`** and logs the total number of products and each item.

## gRPC API

Service definition (simplified):

```proto
service ProductService {
  rpc Fetch(FetchRequest) returns (FetchResponse);
  rpc List(ListRequest) returns (ListResponse);
}
```

- **Fetch**
  - **Request**: `FetchRequest { string url }`
  - **Behavior**:
    - Downloads a `;`-separated CSV.
    - Uses product name as the document ID.
    - Upserts documents in MongoDB and tracks created/updated counts.
  - **Response**: `FetchResponse { int32 created; int32 updated; string error; }`

- **List**
  - **Request**: `ListRequest { PagingParams paging; SortingParams sorting; }`
  - **Response**: `ListResponse { repeated Product products; int32 total; }`

Supported sorting fields: **`name`**, **`price`**.

## Regenerating Protobuf Code

If you modify `api/proto/products.proto`, regenerate Go code with:

```bash
protoc \
  --go_out=. \
  --go-grpc_out=. \
  api/proto/products.proto
```

Make sure the `go_package` option in the proto file matches your module path.

## Project Structure

```text
.
├── api/
│   └── proto/
│       ├── products.proto        # gRPC service and messages
│       └── products.pb.go        # Generated Go code
├── cmd/
│   ├── server/                   # gRPC server entrypoint
│   └── client/                   # gRPC client example
├── internal/
│   └── app/
│       └── server/               # ProductService implementation
├── docker-compose.yml
├── go.mod
└── README.md
```

## Notes

- Products are stored in MongoDB with the product name as `_id`.
- CSV files are expected to use `;` as a delimiter and contain at least `name` and `price` columns.
