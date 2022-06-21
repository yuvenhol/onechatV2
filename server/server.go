package server

import (
	"hash/crc32"
	"log"
	"onechat/protocol/domain"
	"strings"
	"time"

	"github.com/panjf2000/gnet"
)

type OneChatServer struct {
	*gnet.EventServer
}

var sessionMap = make(map[uint32]*Session)

func (ocs *OneChatServer) OnInitComplete(s gnet.Server) (action gnet.Action) {
	log.Println("server started")
	return
}

// 建立连接
func (ocs *OneChatServer) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	log.Println("someone enter")
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
	case domain.COMMAND:
		handleCommand(req, c)
	}
}

//生成session
func setSession(username string, c gnet.Conn) {
	sessionId := crc32.ChecksumIEEE([]byte(username + time.Nanosecond.String()))
	session := &Session{sessionId: sessionId, username: username, con: c}
	sessionMap[sessionId] = session
	c.SetContext(sessionId)
}

//清空session
func clearSession(c gnet.Conn) {
	if c.Context() != nil {
		sessionId := c.Context().(uint32)
		delete(sessionMap, sessionId)
	}

}

//处理聊天
func handleTalk(req *domain.REQ, c gnet.Conn) {
	sessionId := c.Context().(uint32)
	for k, v := range sessionMap {
		if k != sessionId {
			domain.SendAck(v.username, req.Content, v.con)
		}
	}
}

//处理命令
func handleCommand(req *domain.REQ, c gnet.Conn) {
	var message string
	sessionId := c.Context().(uint32)

	switch req.Content {
	case "\\who":
		message = whoOnline(sessionId)
	default:
		message = "暂时不支持该功能"
	}

	domain.SendAck("system", message, c)
}

func whoOnline(sessionId uint32) string {
	oneLineUserNames := make([]string, 0, len(sessionMap))
	for _, v := range sessionMap {
		oneLineUserNames = append(oneLineUserNames, v.username)
	}
	return strings.Join(oneLineUserNames, ",")
}

//用户连接session
type Session struct {
	sessionId uint32
	//用户名
	username string
	//连接
	con gnet.Conn
}
