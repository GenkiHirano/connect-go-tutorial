server:
	go run ./cmd/server/main.go

client:
	go run ./cmd/client/main.go

lint:
	go fmt ./cmd/...
