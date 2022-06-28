package server

import (
	"errors"
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
var historyQueue *ACKQueue

func (ocs *OneChatServer) OnInitComplete(s gnet.Server) (action gnet.Action) {
	historyQueue = NewACKQueue(10)
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
	if req == nil {
		return
	}
	log.Printf("发送内容:%+v", req)
	service(req, c)
	return
}

func (es *OneChatServer) Tick() (delay time.Duration, action gnet.Action) {
	delay = 200 * time.Millisecond
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
	senderSession, ok := sessionMap[sessionId]
	if !ok {
		log.Println("senderSession is null")
		return
	}
	ack := &domain.ACK{Username: senderSession.username, Content: req.Content}
	//broadcast message
	for k, v := range sessionMap {
		if k != sessionId {
			domain.SendAck(ack, v.con)
		}
	}
	//save to history
	historyQueue.Put(*ack)
}

//处理命令
func handleCommand(req *domain.REQ, c gnet.Conn) {
	sessionId := c.Context().(uint32)

	switch req.Content {
	case "\\who":
		message := whoOnline(sessionId)
		ack := &domain.ACK{Username: "system", Content: message}
		domain.SendAck(ack, c)
	case "\\his":
		for _, ack := range historyQueue.GetAcks() {
			domain.SendAck(&ack, c)
			log.Printf("send %v\n", ack)
		}
	default:
		message := "暂时不支持该功能"
		ack := &domain.ACK{Username: "system", Content: message}
		domain.SendAck(ack, c)
	}

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

// ack loop-queue

// This queue is used to store history message.
type ACKQueue struct {
	front int
	tail  int
	size  int
	cap   int
	acks  []domain.ACK
}

func NewACKQueue(cap int) *ACKQueue {
	return &ACKQueue{cap: cap, acks: make([]domain.ACK, cap)}
}

//check the ACKQueue is inited
func (q *ACKQueue) check() error {
	if q.cap == 0 || len(q.acks) == 0 {
		return errors.New("ACKQueue is not inited")
	}
	return nil
}

//enter queue
func (q *ACKQueue) Put(ack domain.ACK) {
	q.check()
	//指针碰撞，弹出start
	if q.tail == q.front && q.size > 0 {
		q.Pop()
	}
	q.size++
	q.acks[q.tail] = ack
	q.tail = (q.tail + 1) % q.cap
}

//out queue
func (q *ACKQueue) Pop() (r domain.ACK) {
	q.check()
	if q.front == q.tail && q.size == 0 {
		return
	}
	r = q.acks[q.front]
	q.front = (q.front + 1) % q.cap
	q.size--
	return
}

func (q *ACKQueue) GetAcks() (result []domain.ACK) {
	start, end := q.front, q.tail
	if q.size == 0 && start == end {
		return
	}
	if start < end {
		result = make([]domain.ACK, end-start)
		copy(result, q.acks[start:end])
	} else if start >= end {
		result = make([]domain.ACK, 0, end+q.cap-start)
		result = append(result, q.acks[start:]...)
		result = append(result, q.acks[:end]...)
	}
	return
}
