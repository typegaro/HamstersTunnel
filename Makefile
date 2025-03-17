GO = go

# Install dependencies (tidy up the go.mod file)
install:
	$(GO) mod tidy

# Run the server
server: build-server
	./cmd/server/bin/server

cli: build-cli
	./cmd/cli/bin/cli
	
daemon: build-daemon
	./cmd/daemon/bin/daemon
# Run all tests
test:
	$(GO) test ./... -v

build-server:
	$(GO) build -o ./cmd/server/bin/server ./cmd/server/main.go

build-daemon:
	$(GO) build -o ./cmd/daemon/bin/daemon ./cmd/daemon/main.go 

build-cli:
	$(GO) build -o ./cmd/cli/bin/cli ./cmd/cli/main.go 
# Build the server executable
build: install build-daemon build-cli build-server

# Clean the build (remove the compiled executable)
clean:
	rm -f ./cmd/server/bin/server

# Clean the Go mod dependencies (optional, useful for cleaning up unnecessary dependencies)
clean-mod:
	$(GO) clean -modcache

