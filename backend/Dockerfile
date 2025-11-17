# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy backend source code
COPY backend/ ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates python3 py3-pip

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy Python ingestion scripts
COPY backend/ingestion/ ./ingestion/

# Install Python dependencies
RUN pip3 install --break-system-packages -r ./ingestion/requirements.txt

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
