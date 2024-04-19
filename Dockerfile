# Use the official Go image to create a build container
FROM golang:1.17 AS builder

WORKDIR /workspaces/listaway

# Copy the Go module files
COPY go.mod .
COPY go.sum .

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/listaway ./cmd/listaway

# Use a minimal base image to run the application
FROM alpine:latest

WORKDIR /root/

# Copy the built executable from the previous stage
COPY --from=builder /workspaces/listaway/bin/listaway .

# Copy over sql files
COPY --from=builder /workspaces/listaway/db/* ./db/

EXPOSE 8080

# Run the Go application
CMD ["./listaway"]
