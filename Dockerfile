FROM golang:1.23 AS builder

# Install git
RUN apt-get update && apt-get install -y git

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build the application with the version flag
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$(git describe --tags --always --dirty --abbrev=8)" -a -installsuffix cgo -o kubectl-ranching.farm .