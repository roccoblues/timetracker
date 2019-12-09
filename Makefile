GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOTIDY=$(GOCMD) mod tidy
BINARY_NAME=tt

.PHONY: all
all: tidy test build

GIT_COMMIT := $(shell git rev-list -1 HEAD)
BUILD_ARGS=-ldflags "-X main.version=$(GIT_COMMIT)"

.PHONY: build
build:
	$(GOBUILD) $(BUILD_ARGS) -o $(BINARY_NAME) -v cmd/tt/*

.PHONY: test
test:
	$(GOTEST) -v -race ./...

.PHONY: clean
clean:
	rm -f $(BINARY_NAME)

.PHONY: tidy
tidy:
	$(GOTIDY)
