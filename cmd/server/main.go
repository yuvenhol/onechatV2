package main

import (
	"log"
	"onechat/server"

	"github.com/panjf2000/gnet"
)

func main() {
	onechat := new(server.OneChatServer)
	log.Fatal(gnet.Serve(onechat, "tcp://127.0.0.1:9002", gnet.WithMulticore(false)))
}
