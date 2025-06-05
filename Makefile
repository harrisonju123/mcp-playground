GO_BIN=mcpxd

.PHONY: dev test lint docker

dev:                                     ## Local dev loop
	air -c .air.toml & pnpm --prefix web dev

test:                                    ## Run all tests
	go test ./...
	pnpm --prefix web test --if-present

lint:
	golangci-lint run ./...
	pnpm --prefix web lint --fix

docker:
	docker build -f Dockerfile.daemon -t $(GO_BIN):dev .

.DEFAULT_GOAL := dev