# Build stage: use Alpine to produce a statically linked binary
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o /mcpxd ./cmd/mcpxd

# Final stage: use scratch to run the static binary
FROM scratch

COPY --from=builder /mcpxd /mcpxd
COPY --from=builder /app/config /config
ENTRYPOINT ["/mcpxd"]