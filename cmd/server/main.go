package main

import (
	"fmt"
	"log"
	"onechat/server"

	"github.com/panjf2000/gnet"
)

func main() {
	defer func() {
		if err := recover(); err != nil {

			fmt.Println("catch error ", err)
		}
	}()
	onechat := new(server.OneChatServer)
	log.Fatal(gnet.Serve(onechat, "tcp://:5211", gnet.WithMulticore(false)))
}
