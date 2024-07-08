# Use the official Go image to create a build container
FROM golang:1.22-bookworm AS builder

WORKDIR /workspaces/listaway

# Copy the Go module files
COPY go.mod .
COPY go.sum .

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Install npm
RUN apt-get update && apt-get install -y curl \
    && curl -fsSL https://deb.nodesource.com/setup_20.x | bash - \
    && apt-get install -y nodejs

# Install yarn
RUN npm install --global yarn

# Install npm dependencies
RUN yarn --cwd web install

# Build the static assets using npm
RUN yarn --cwd web run build-prod

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
