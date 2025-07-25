# Concurrency Golang FullCycle Lab 03 - Auction System

This project implements an auction system in Go with concurrency features, using MongoDB as the database and Gin as the web framework.

## Prerequisites

- **Go 1.24.3** or higher
- **Docker** and **Docker Compose**
- **MongoDB** (if running locally without Docker)

## Development Environment Setup

### 1. Clone the repository

```bash
git clone <repository-url>
cd concurrency-golang-fullcycle-lab03
```

### 2. Environment Variables Configuration

The project uses an `.env` file located at `cmd/auction/.env`. The default environment variables are:

```env
MONGODB_URL=mongodb://localhost:27017
MONGODB_DB=auctionsDB
BATCH_INSERT_INTERVAL=7m
MAX_BATCH_SIZE=10
AUCTION_INTERNAL=5m
```

**Important:** The `.env` file is already configured with default values for local development.

### 3. Install Dependencies

```bash
go mod download
```

### 🐳 Running with Docker (Recommended)

### Option 1: Docker Compose (Simplest)

Run the complete system (application + MongoDB) with a single command:

```bash
docker-compose up --build
```

This will:
- Build the application image
- Start MongoDB on port 27017
- Start the application on port 8080

### Option 2: MongoDB only via Docker

If you prefer to run only MongoDB via Docker and the application locally:

```bash
# Start MongoDB only
docker run -d \
  --name mongodb \
  -p 27017:27017 \
  -e MONGO_INITDB_DATABASE=auctionsDB \
  mongo:7.0
```

### Running Locally (Without Docker)

### 1. Local MongoDB

Ensure MongoDB is running locally on port 27017:

```bash
# On macOS with Homebrew
brew services start mongodb-community

# On Ubuntu/Debian
sudo systemctl start mongod

# On Windows
net start MongoDB
```

### 2. Run the Application

```bash
# Navigate to the application directory
cd cmd/auction

# Run the application
go run main.go
```

Or compile and run:

```bash
# Compile
go build -o auction cmd/auction/main.go

# Run
./auction
```

### API Endpoints

The application will be available at `http://localhost:8080` with the following endpoints:

### Auctions
- `GET /auctions` - List all auctions
- `GET /auctions/:auctionId` - Find auction by ID
- `POST /auctions` - Create new auction
- `GET /auction/winner/:auctionId` - Find winning bid by auction ID

### Bids
- `GET /bid/:auctionId` - Find bids by auction ID
- `POST /bid` - Create new bid

### Users
- `GET /user/:userId` - Find user by ID

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with detailed coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```


