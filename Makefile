GOCACHE ?= /tmp/go-build

.PHONY: generate validate-manifest docs docs-validate test test-acceptance coverage-report

generate:
	GOCACHE=$(GOCACHE) go run ./cmd/awxgen generate

docs:
	GOCACHE=$(GOCACHE) go run ./cmd/awxgen docs

docs-validate:
	GOCACHE=$(GOCACHE) go run ./cmd/awxgen docs-validate

validate-manifest:
	GOCACHE=$(GOCACHE) go run ./cmd/awxgen validate

coverage-report:
	GOCACHE=$(GOCACHE) go run ./cmd/awxgen report

test:
	GOCACHE=$(GOCACHE) go test ./...

test-acceptance:
	@set -a; \
	if [ -f .env ]; then . ./.env; fi; \
	set +a; \
	GOCACHE=$(GOCACHE) go test ./internal/acceptance -run TestAcceptance -v
