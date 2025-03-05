GO = go

# Install dependencies (tidy up the go.mod file)
install:
	$(GO) mod tidy

# Run the server
server: build-server
	./cmd/server/bin/server

client: build-client
	./cmd/client/bin/client

# Run all tests
test:
	$(GO) test ./... -v

build-server:
	$(GO) build -o ./cmd/server/bin/server ./cmd/server/main.go

build-client:
	$(GO) build -o ./cmd/client/bin/client ./cmd/client/main.go 

# Build the server executable
build: install
	$(GO) build -o ./cmd/server/bin/server ./cmd/server/main.go
	$(GO) build -o ./cmd/client/bin/client ./cmd/client/main.go 

# Clean the build (remove the compiled executable)
clean:
	rm -f ./cmd/server/bin/server

# Clean the Go mod dependencies (optional, useful for cleaning up unnecessary dependencies)
clean-mod:
	$(GO) clean -modcache

