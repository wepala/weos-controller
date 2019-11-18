# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=weos-controller

all: deps test
build:
		$(GOCMD) run -ldflags "-X main.version=dev -X main.build=12082019" main.go
		$(GOBUILD) -v -o $(BINARY_NAME) -ldflags "-X main.version=dev -X main.build=12082019" main.go
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
deps:
		$(GOGET) github.com/golang/dep/cmd/dep
		dep ensure
