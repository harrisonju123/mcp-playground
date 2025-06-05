FROM golang:1.24 as builder
WORKDIR /src
COPY . .
RUN go build -o /bin/mcpxd ./cmd/mcpxd

FROM gcr.io/distroless/base-debian11
COPY --from=builder /bin/mcpxd /mcpxd
EXPOSE 50051
ENTRYPOINT ["/mcpxd"]