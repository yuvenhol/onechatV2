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
	codec := &gnet.LineBasedFrameCodec{}
	log.Fatal(gnet.Serve(onechat, "tcp://:5211", gnet.WithCodec(codec), gnet.WithMulticore(false)))
}
