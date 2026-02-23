GOCACHE ?= /tmp/go-build

.PHONY: generate validate-manifest docs docs-validate docs-verify-online docs-validate-quality test test-acceptance coverage-report build

QUALITY_SUMMARY ?=
MAX_QUALITY_PASSES ?= 3

generate:
	GOCACHE=$(GOCACHE) go run ./cmd/awxgen generate

docs:
	GOCACHE=$(GOCACHE) go run ./cmd/awxgen docs

docs-validate:
	GOCACHE=$(GOCACHE) go run ./cmd/awxgen docs-validate

docs-verify-online:
	GOCACHE=$(GOCACHE) go run ./cmd/awxgen docs-verify-online

docs-validate-quality:
	@if [ -z "$(QUALITY_SUMMARY)" ]; then echo "QUALITY_SUMMARY is required (for example openspec/changes/<change>/implementation-summary.md)"; exit 1; fi
	GOCACHE=$(GOCACHE) go run ./cmd/awxgen docs-validate-quality --summary "$(QUALITY_SUMMARY)" --max-passes "$(MAX_QUALITY_PASSES)"

validate-manifest:
	GOCACHE=$(GOCACHE) go run ./cmd/awxgen validate

coverage-report:
	GOCACHE=$(GOCACHE) go run ./cmd/awxgen report

build:
	mkdir -p dist
	GOCACHE=$(GOCACHE) go build -o dist/terraform-provider-awx-iwd ./cmd/terraform-provider-awx-iwd

test:
	GOCACHE=$(GOCACHE) go test ./...

test-acceptance:
	@set -a; \
	if [ -f .env ]; then . ./.env; fi; \
	set +a; \
	TF_ACC=1 GOCACHE=$(GOCACHE) go test ./internal/acceptance ./internal/provider -run TestAcceptance -v
