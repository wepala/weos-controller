# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=weos-controller

all: deps test build
build:
		$(GOBUILD) -o $(BINARY_NAME) -v
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
deps:
		$(GOGET) github.com/golang/dep/cmd/dep
		dep ensure
