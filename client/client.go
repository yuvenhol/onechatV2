package client

import (
	"log"
	"onechat/protocol/domain"

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
