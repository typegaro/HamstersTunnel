
# Makefile per gestire il progetto Go

# Variabili
GO = go

# Comando per avviare il server
server:
	$(GO) run cmd/server/main.go

# Comando per eseguire i test in tutte le cartelle
test:
	$(GO) test ./... -v

# Comando per compilare il progetto
build:
	$(GO) build -o ./cmd/server/bin/server ./cmd/server/main.go

# Pulizia dei file generati
clean:
	rm ./cmd/server/bin/server
