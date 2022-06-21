package server

import (
	"hash/crc32"
	"log"
	"onechat/protocol/domain"
	"time"

	"github.com/panjf2000/gnet"
)

type OneChatServer struct {
	*gnet.EventServer
}

var sessionMap = make(map[uint32]*Session)

// 建立连接
func (ocs *OneChatServer) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	log.Println("enter")
	return
}

// 断开连接
//清除session
func (ocs *OneChatServer) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	clearSession(c)
	return
}

//正常数据传输
//给其他在线用户发消息
func (ocs *OneChatServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	req := domain.Parse2Req(frame)
	log.Printf("发送内容:%+v", req)
	service(req, c)
	return

}

//处理请求
func service(req *domain.REQ, c gnet.Conn) {
	switch req.Action {
	case domain.ENTER:
		setSession(req.Content, c)
	case domain.TALK:
		handleTalk(req, c)
	case domain.QUERY:
		//TODO: \who \history
	}
}

//生成session
func setSession(username string, c gnet.Conn) {
	sessionId := crc32.ChecksumIEEE([]byte(username + time.Nanosecond.String()))
	session := &Session{sessionId: sessionId, username: username, con: c}
	sessionMap[sessionId] = session
	c.SetContext(sessionId)
}

func clearSession(c gnet.Conn) {
	if c.Context() != nil {
		token := c.Context().(uint32)
		delete(sessionMap, token)
	}

}

func handleTalk(req *domain.REQ, c gnet.Conn) {
	token := c.Context().(uint32)
	for k, v := range sessionMap {
		if k != token {
			domain.SendAck(v.username, req.Content, v.con)
		}
	}
}

//用户连接session
type Session struct {
	sessionId uint32
	//用户名
	username string
	//连接
	con gnet.Conn
}
