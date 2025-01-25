# Stage 1: Build the application
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o main .

# Stage 2: Create a minimal runtime image
FROM gcr.io/distroless/base-debian12

# Set environment variables
ENV PORT=8100

# Set the working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Expose the application's port
EXPOSE 8100

# Command to run the application
CMD ["./main"]
