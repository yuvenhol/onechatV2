package domain

import (
	"encoding/json"
	"log"
	"onechat/protocol/security"
	"time"

	"github.com/panjf2000/gnet"
)

//请求
//client -> server 消息体
type REQ struct {
	//动作
	Action ReqAction
	//请求内容
	Content string
}

//回应
//server -> client 消息体
type ACK struct {
	//用户名
	Username string
	//请求内容
	Content string
}

//发回应
func SendAck(username string, content string, c gnet.Conn) {
	if len(content) == 0 {
		log.Println("发送内容为空")
	}
	ack := ACK{Username: username, Content: content}
	bs, err := json.Marshal(ack)
	if err != nil {
		panic(err)
	}
	//加密发送
	err = c.AsyncWrite(security.Encrypt(bs))
	if err != nil {
		panic(err)
	}
}

//发送请求
func SendReq(action ReqAction, content string, c gnet.Conn) {
	if len(content) == 0 {
		log.Println("发送内容为空")
	}
	req := REQ{Action: action, Content: content}
	bs, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}
	//加密发送
	err = c.SendTo(security.Encrypt(bs))
	if err != nil {
		panic(err)
	}
	time.Sleep(100 * time.Millisecond)
}

//解析请求
func Parse2Req(frame []byte) *REQ {
	//解密
	req := &REQ{}
	b := security.Decrypt(frame)
	if json.Unmarshal(b, req) != nil {
		log.Println("json解析异常")
	}
	return req
}

//解析回应
func Parse2ACK(frame []byte) *ACK {
	ack := &ACK{}
	if json.Unmarshal(security.Decrypt(frame), ack) != nil {
		log.Println("json解析异常")
	}
	return ack
}

//请求动作类型枚举
type ReqAction int

const (
	//进入聊天室
	ENTER ReqAction = iota
	//聊天
	TALK
	//命令
	COMMAND
)
