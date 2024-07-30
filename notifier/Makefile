GOTOOLCHAIN := $(shell which go)

.PHONY: build
build:
	@echo Using Go at: $(GOTOOLCHAIN)
	@echo "Building the project for ${GOOS}/${GOARCH}..."
	GOOS=${GOOS} GOARCH=${GOARCH} $(GOTOOLCHAIN) build -o build/notifier cmd/notifier/main.go

.PHONY: coverage
coverage:
	@echo "Running tests with race detection and coverage reporting..."
	@FILES=$$(go list ./... | grep -v /mock | grep -v /test); if [ -z "$$FILES" ]; then \
		echo "No Go files to test after filtering"; \
	else \
		go test -race -covermode=atomic -coverprofile=./coverage/coverage.out -count=5 $$FILES; \
		go tool cover -func=./coverage/coverage.out; \
		echo "Generating HTML coverage report..."; \
		go tool cover -html=./coverage/coverage.out -o ./coverage/coverage.html; \
	fi
	@echo "Done!"

.PHONY: lint
lint:
	@echo "Running golangci-lint..."
	@golangci-lint run
	@echo "Done!"
