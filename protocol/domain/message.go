package domain

import (
	"encoding/json"
	"log"
	"onechat/protocol/security"

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
func SendAck(ack *ACK, c gnet.Conn) {
	if len(ack.Content) == 0 {
		log.Println("发送内容为空")
		return
	}
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

func SendAckV(acks []ACK, c gnet.Conn) {
	qrBs := make([][]byte, 0)
	for _, ack := range acks {
		bs, err := json.Marshal(ack)
		if err != nil {
			panic(err)
		}
		qrBs = append(qrBs, security.Encrypt(bs))
	}
	c.AsyncWritev(qrBs)

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
	err = c.AsyncWrite(security.Encrypt(bs))
	if err != nil {
		panic(err)
	}
}

//解析请求
func Parse2Req(frame []byte) *REQ {
	//解密
	req := &REQ{}
	b := security.Decrypt(frame)
	if err := json.Unmarshal(b, req); err != nil {
		log.Printf("json解析异常%s", err.Error())
	}
	return req
}

//解析回应
func Parse2ACK(frame []byte) *ACK {
	ack := &ACK{}
	b := security.Decrypt(frame)
	if err := json.Unmarshal(b, ack); err != nil {
		log.Printf("json解析异常%s\nbytes:%s", err.Error(), string(b))
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
