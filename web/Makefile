# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=tizizz-web
WX_VERIFY=MP_verify_VmrbuqOYQnB6tzAS.txt
# Build the project
all: build

build:
	# build frontend
	# cd ui && pnpm build
	# cd ..
	cp -f ui/$(WX_VERIFY) ui/out/$(WX_VERIFY)
	$(GOBUILD) -o $(BINARY_NAME) -v

# Run tests
test:
	$(GOTEST) -v ./...

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Install dependencies
deps:
	$(GOGET) -v ./...

# Format the code
fmt:
	$(GOCMD) fmt ./...

.PHONY: all build clean test deps fmt