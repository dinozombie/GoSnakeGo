# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
BINARY_NAME=$(GOPATH)/bin/gosnakego

all: install
build: 
	$(GOBUILD) -o $(BINARY_NAME) -v
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
run:
	$(GOBUILD) -o $(BINARY_NAME) -v
	$(BINARY_NAME)
deps:
	$(GOGET) github.com/rthornton128/goncurses
install: 
	$(GOGET) github.com/rthornton128/goncurses
	$(GOBUILD) -o $(BINARY_NAME) -v

