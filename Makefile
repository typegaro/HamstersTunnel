GO = go

server:
	$(GO) run cmd/server/main.go

test:
	$(GO) test ./... -v

build:
	$(GO) build -o ./cmd/server/bin/server ./cmd/server/main.go

clean:
	rm ./cmd/server/bin/server
