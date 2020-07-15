# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=weos-controller

all: build
build:
		$(GOCMD) run -ldflags "-X main.version=dev -X main.build=11232019" main.go
		$(GOBUILD) -v -o $(BINARY_NAME) -ldflags "-X main.version=dev -X main.build=11232019" main.go
		chmod u+x $(BINARY_NAME)
		chmod o+x $(BINARY_NAME)
test:
		$(GOCMD) fmt
		$(GOCMD) vet -v
		$(GOTEST) -v ./...
clean:
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
run:
		$(GOBUILD) -o $(BINARY_NAME) -v
		./$(BINARY_NAME)