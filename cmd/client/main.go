package main

import (
	"bufio"
	"flag"
	"log"
	"onechat/client"
	"onechat/protocol/domain"
	"os"
	"strings"

	"github.com/panjf2000/gnet"
)

func main() {

	var address = flag.String("add", "127.0.0.1:9002", "TCP listening address")
	var username = flag.String("u", "default", "username")
	flag.Parse()

	client, err := gnet.NewClient(new(client.OneChatClient))
	if err != nil {
		log.Fatal("连接异常")
		return
	}
	client.Start()
	//tcp拨号
	conn, err := client.Dial("tcp", *address)
	if err != nil {
		log.Fatal("连接异常")
		return
	}
	defer conn.Close()
	log.Println("wecome to onechat :) ")
	//登录
	domain.SendReq(domain.ENTER, *username, conn)
	log.Println("logined")
	//先查询一次在线
	go domain.SendReq(domain.COMMAND, "\\who", conn)
	//控制台交互
	ReadConsole(conn)
}

//控制台读取键盘输入
func ReadConsole(con gnet.Conn) {

	for {
		// 从标准输入读取字符串，以\n为分割
		text, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			break
		}
		//滤掉空数据
		text = strings.TrimSpace(text)
		if len(text) == 0 {
			continue
		}
		//默认对话行为
		action := domain.TALK
		//命令行为
		if strings.HasPrefix(text, "\\") {
			action = domain.COMMAND
		}
		domain.SendReq(action, text, con)
	}

}
