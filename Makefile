# ---------------------------------------------------------------------
# Makefile – workspace-friendly dev helpers for mcp-agent-poc
# ---------------------------------------------------------------------

# 1) List every Go module directory that should be vetted / tested.
GO_MODULES := .

# 2) Name of the daemon binary and Docker image tag
GO_BIN      := mcpxd
IMAGE_TAG   := $(GO_BIN):dev

# ---------------------------------------------------------------------
# Targets
# ---------------------------------------------------------------------

.PHONY: dev gotest govet lint docker clean gen

## dev – run backend and frontend in watch-mode
dev:
	# Backend hot-reload (uncomment Air if you use it)
	# air &
	go run ./cmd/mcpxd &                      \
	pnpm --prefix web dev

## gotest – run go test in every module
gotest:
	@for m in $(GO_MODULES); do \
		echo "› go test $$m/..."; \
		go test $$m/...; \
	done

## govet – run go vet in every module
govet:
	@for m in $(GO_MODULES); do \
		echo "› go vet $$m/..."; \
		go vet $$m/...; \
	done

## lint – golangci-lint + eslint
lint:
	golangci-lint run ./...            # Go
	pnpm --prefix web lint --fix       # Web

## docker – build distroless image containing the daemon
docker:
	docker build -f Dockerfile.daemon -t $(IMAGE_TAG) .

## clean – tidy modules & remove Docker image
clean:
	@for m in $(GO_MODULES); do \
		( cd $$m && go mod tidy ); \
	done
	docker rmi -f $(IMAGE_TAG) || true

gen:
	@echo "🔄  Running go generate for all proto‐generate directives…"
	protoc \
      --proto_path=api \
      --go_out=api/gen --go_opt=paths=source_relative \
      --go-grpc_out=api/gen --go-grpc_opt=paths=source_relative \
      api/aggregator.proto

e2e:
	docker compose -f infra/docker-compose.yml up --build -d
	sleep 5             # warm-up
	k6 run test/e2e/k6.js
	./scripts/failover.sh
	k6 run test/e2e/k6.js
	docker compose -f infra/docker-compose.yml down -v

.DEFAULT_GOAL := dev