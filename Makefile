build:
	go build -o cmd/server/server_v2.0.1 cmd/server/main.go
	go build -o cmd/client/client_v2.0.1 cmd/client/main.go


run:
	go run cmd/server/main.go