package server_test

import (
	"fmt"
	"onechat/protocol/domain"
	"onechat/server"
	"testing"
)

func TestLoopQueue(t *testing.T) {
	queue := server.NewACKQueue(10)
	for i := 0; i < 8; i++ {
		ack := &domain.ACK{Username: "aaa", Content: fmt.Sprint(i)}
		queue.Put(*ack)
	}
	fmt.Printf("%v\n", queue.GetAcks())
	for i := 1; i < 6; i++ {
		ack := &domain.ACK{Username: "aaa", Content: fmt.Sprint(i)}
		queue.Put(*ack)
	}
	fmt.Printf("%v\n", queue.GetAcks())
	for i := 1; i < 6; i++ {
		queue.Pop()
	}
	fmt.Printf("%v\n", queue.GetAcks())
}
