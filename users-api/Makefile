.PHONY: test test-integration test-unit test-clean

# Default target
all: test

# Run all tests
test: test-unit test-integration

# Run unit tests
test-unit:
	go test -v ./services/... ./repositories/...

# Run integration tests with Docker
test-integration:
	docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit

# Clean up test containers and volumes
test-clean:
	docker-compose -f docker-compose.test.yml down -v 