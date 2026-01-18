# Subscription Service

A subscription service built with Go and Gin.

## Installation

```bash
git clone <repository-url>
cd subscription-service
```

### Running with Docker Compose

Build and run the service using Docker Compose:

```bash
docker build --tag subscription-service .
```

```bash
docker-compose up --build
```

To stop the service:

```bash
docker-compose down
```

### Running with Docker

Alternatively, run using Docker directly:

```bash
docker build --tag subscription-service .
docker run -p 8080:8080 subscription-service
```

### Running Locally

Initialize dependencies:

```bash
go mod download
```

Run the service:

```bash
go run ./cmd/subscription-service/
```

Or build and run:

```bash
go build -o . ./cmd/subscription-service/
./subscription-service
```

## API Endpoints

Service provides REST API service 

Full documentation in Swagger:
```
http://localhost:8080/swagger/index.html
```

## Configuration

Copy the example environment file and configure your settings:

Edit `.env` with your configuration values.