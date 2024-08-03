# Build stage
FROM golang:1.20-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./
COPY go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app with optimizations for production
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /chat_app ./cmd/main.go

# Final stage
FROM alpine:latest



# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the pre-built binary from the build stage
COPY --from=builder /chat_app .

# Copy .env file into the container
COPY .env /root/.env

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["./chat_app"]