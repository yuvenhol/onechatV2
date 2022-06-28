package client

import (
	"log"
	"onechat/protocol/domain"
	"time"

	"github.com/panjf2000/gnet"
)

type OneChatClient struct {
	*gnet.EventServer
}

func (occ *OneChatClient) React(content []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	a := domain.Parse2ACK(content)
	log.Println(a.Username + "->" + a.Content)
	return
}

func (es *OneChatClient) Tick() (delay time.Duration, action gnet.Action) {
	delay = 200 * time.Millisecond
	return
}

func (es *OneChatClient) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	log.Fatalln("connect is closed")
	return
}
