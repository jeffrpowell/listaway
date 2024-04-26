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

# Compile css file from TailwindCSS classes
RUN curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/download/v3.4.3/tailwindcss-linux-x64
RUN chmod +x tailwindcss-linux-x64
RUN mv tailwindcss-linux-x64 /usr/local/bin/tailwindcss
RUN tailwindcss -i ./web/root.css -o ./internal/handlers/assets/root.css

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/listaway ./cmd/listaway

# Use a minimal base image to run the application
FROM alpine:latest

WORKDIR /root/

# Copy the built executable from the previous stage
COPY --from=builder /workspaces/listaway/bin/listaway .

EXPOSE 8080

# Run the Go application
CMD ["./listaway"]
