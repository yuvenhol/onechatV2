VERSION= v2.0.2
build:
	go build -o cmd/server/server_$(VERSION) cmd/server/main.go
	go build -o cmd/client/client_$(VERSION) cmd/client/main.go
	cp cmd/client/client_$(VERSION) /usr/local/bin/oneChat_client
buildamd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cmd/server/onechat_server_$(VERSION)_linux_amd64 cmd/server/main.go
run:
	go run cmd/server/main.go