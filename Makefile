.PHONY: docker-compose-up unit-tests integration-tests tests run

GO ?= go

# List all packages except the integration test package under ./test
UNIT_PACKAGES := $(shell $(GO) list ./... | grep -v '/test$$')

docker-compose-up:
	@echo "Starting docker compose stack..."
	@docker compose -f deployments/docker-compose.yml up -d 2>/dev/null \
		|| docker-compose -f deployments/docker-compose.yml up -d

unit-tests:
	@echo "Running unit tests..."
	@$(GO) test -v $(UNIT_PACKAGES)

integration-tests:
	@echo "Running integration tests..."
	@$(GO) test -v ./test -count=1

tests: unit-tests integration-tests

run:
	@echo "Running application..."
	@$(GO) run ./cmd

