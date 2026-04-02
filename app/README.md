# GoLedger Challenge - Besu Application

gRPC service that integrates Hyperledger Besu blockchain with PostgreSQL, providing read/write operations on a smart contract with database synchronization.

## Architecture

Clean architecture with three layers:

```
Transport (gRPC) → Service → Repository + Blockchain
```

- **Transport Layer** ([`internal/transport/grpc`](internal/transport/grpc)): gRPC server handling requests
- **Service Layer** ([`internal/service`](internal/service)): Business logic coordinating blockchain and database
- **Data Layer**: 
  - **Blockchain** ([`internal/blockchain`](internal/blockchain)): Source of truth - smart contract operations
  - **Repository** ([`internal/repository`](internal/repository)): PostgreSQL cache for fast reads

### Design Decisions

1. **Blockchain as Source of Truth**: All writes go to blockchain first, database is best-effort cache
2. **Transaction Hash Return**: `SetValue` returns blockchain tx hash for verifiability
3. **Graceful Error Handling**: Database failures on cache updates are logged but don't fail operations
4. **Configuration Package**: Centralized config with environment variable support
5. **No External Dependencies**: Uses standard library + `go-ethereum` for idiomatic Go

## Prerequisites

- Go 1.21+
- Docker & Docker Compose
- Running Besu network (see [root README](../README.md))
- Deployed SimpleStorage contract

## Quick Start

### 1. Start Infrastructure

```bash
# From repository root
make devnet-deploy  # Starts Besu network and deploys contract

# Start PostgreSQL
docker-compose up -d
```

### 2. Configure Application

Copy and edit environment variables:

```bash
cd app
cp .env.example .env
# Edit CONTRACT_ADDRESS with deployed contract address from make devnet-deploy
```

### 3. Run Application

```bash
# Install dependencies
go mod download

# Run database migrations
go run migrations/*.go

# Start gRPC server
go run cmd/server/main.go
```

Expected output:
```
✅ Connected to PostgreSQL
✅ Connected to Besu (Chain ID: 1337, Block: 2616)
✅ Contract loaded at 0x42699a7612a82f1d9c36148af9c77354759b210b
✅ Service initialized with blockchain integration
🚀 gRPC server listening on :50051
```

## API Reference

The service exposes 4 gRPC methods (reflection enabled).

### Testing with Postman

1. **Create new gRPC Request**
   - Click **New** → **gRPC Request**
   - Enter server URL: `localhost:50051`
   - **Use Server Reflection** is enabled by default (discovers methods automatically)

2. **Select Method**
   - Choose `storage.StorageService/SetValue` (or any other method)

3. **Send Request**
   - For `SetValue`, enter JSON message:
     ```json
     {
       "value": 42
     }
     ```
   - Click **Invoke**

4. **View Response**
   - See transaction hash, values, or sync status

### Testing with Insomnia

1. **New Request** → Select **gRPC**
2. **URL**: `localhost:50051`
3. Insomnia auto-discovers methods via reflection
4. Select method and enter JSON body
5. **Send**

### Testing with grpcurl (CLI)

#### SetValue
Write a value to the smart contract.

```bash
grpcurl -plaintext -d '{"value": 42}' localhost:50051 storage.StorageService/SetValue
```

Response:
```json
{
  "txHash": "0x123abc..."
}
```

#### GetValue
Read current value from blockchain.

```bash
grpcurl -plaintext localhost:50051 storage.StorageService/GetValue
```

Response:
```json
{
  "value": 42
}
```

#### SyncValue
Synchronize blockchain value to database.

```bash
grpcurl -plaintext localhost:50051 storage.StorageService/SyncValue
```

Response:
```json
{
  "success": true,
  "message": "database cache synchronized with blockchain"
}
```

#### CheckValue
Compare blockchain value with database cache.

```bash
grpcurl -plaintext -d '{"value": 42}' localhost:50051 storage.StorageService/CheckValue
```

Response:
```json
{
  "inSync": true
}
```

## Configuration

All settings via environment variables (see [`.env.example`](.env.example)):

| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_URL` | `postgres://admin:admin123@localhost:5433/challenge_besu?sslmode=disable` | PostgreSQL connection string |
| `BESU_RPC_URL` | `http://localhost:8545` | Besu JSON-RPC endpoint |
| `CONTRACT_ADDRESS` | `0x42699a7612a82f1d9c36148af9c77354759b210b` | SimpleStorage contract address |
| `PRIVATE_KEY` | `0x8f2a5594903...` | Ethereum private key (pre-funded) |
| `GRPC_PORT` | `50051` | gRPC server port |

## Development

### Project Structure

```
app/
├── cmd/server/          # Application entry point
├── internal/
│   ├── blockchain/      # Besu client and contract bindings
│   ├── config/          # Configuration package
│   ├── database/        # PostgreSQL connection pool
│   ├── domain/          # Domain models
│   ├── repository/      # Database operations
│   ├── service/         # Business logic
│   └── transport/grpc/  # gRPC server implementation
├── proto/               # Protocol buffer definitions
├── gen/pb/              # Generated protobuf code
└── migrations/          # Database schema
```

### Running Tests

```bash
go test ./internal/service -v
go test ./internal/repository -v
```

### Regenerating Protobuf

```bash
protoc --go_out=. --go-grpc_out=. proto/storage.proto
```

## How It Works

1. **Write Flow** (`SetValue`):
   - Sends transaction to blockchain smart contract
   - Waits for transaction to be mined
   - Caches value in PostgreSQL
   - Returns transaction hash

2. **Read Flow** (`GetValue`):
   - Reads from blockchain (source of truth)
   - Updates database cache
   - Returns value

3. **Sync Flow** (`SyncValue`):
   - Explicitly synchronizes blockchain → database
   - Useful after external contract interactions

4. **Check Flow** (`CheckValue`):
   - Compares blockchain value with database
   - Returns sync status

## Troubleshooting

**Server won't start:**
- Verify Besu network is running: `docker logs -f besu.node-1`
- Check PostgreSQL: `docker-compose ps`
- Validate contract address in `.env`

**Transaction fails:**
- Ensure private key has funds
- Check Besu node is mining blocks
- Verify contract is deployed correctly

**Database out of sync:**
- Run `SyncValue` endpoint
- Check logs for cache update errors
