build:
	go build -o cmd/server/server_v2.0.1 cmd/server/main.go
	go build -o cmd/client/client_v2.0.1 cmd/client/main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cmd/server/server_v2.0.1_linux_x64 cmd/server/main.go
run:
	go run cmd/server/main.go