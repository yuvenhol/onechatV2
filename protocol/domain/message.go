package domain

import (
	"encoding/json"
	"log"

	"github.com/panjf2000/gnet"
)

type REQ struct {
	//动作
	Action ReqAction
	//请求内容
	Content string
}

type ACK struct {
	//用户名
	Username string
	//请求内容
	Content string
}

type ReqAction int

func SendAck(username string, content string, c gnet.Conn) {
	ack := ACK{Username: username, Content: content}
	bs, err := json.Marshal(ack)
	if err != nil {
		panic(err)
	}
	err = c.AsyncWrite(bs)
	if err != nil {
		panic(err)
	}
}

func SendReq(action ReqAction, content string, c gnet.Conn) {
	req := &REQ{Action: action, Content: content}
	bs, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}
	err = c.AsyncWrite(bs)
	if err != nil {
		panic(err)
	}
}

func Parse2Req(frame []byte) *REQ {
	var req REQ
	if json.Unmarshal(frame, &req) != nil {
		log.Fatalln("json解析异常")
	}
	return &req
}

func Parse2ACK(frame []byte) *ACK {
	var ack ACK
	if json.Unmarshal(frame, &ack) != nil {
		log.Fatalln("json解析异常")
	}
	return &ack
}

const (
	//进入聊天室
	ENTER ReqAction = iota
	//聊天
	TALK
	//查询
	QUERY
)
