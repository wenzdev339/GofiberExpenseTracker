# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies via proxy
RUN GOPROXY=https://goproxy.io,direct GONOSUMDB="*" go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/server ./cmd/main.go

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Copy binary and migrations
COPY --from=builder /app/server .
COPY --from=builder /app/migrations ./migrations

# Create non-root user
RUN adduser -D -g '' appuser
USER appuser

EXPOSE 3000

ENTRYPOINT ["./server"]
