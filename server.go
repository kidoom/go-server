package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip string
	Post int

	//online map for user
	OnlineMap map[string]*User
	mapLock sync.RWMutex

	//消息广播 channel
	Message chan string
}

// NewServer 创建server接口
func NewServer(ip string,post int) *Server{
	server := &Server{
		Ip: ip,
		Post: post,
		OnlineMap: make(map[string]*User),
		Message: make(chan string),
	}
	return  server
}

// ListenMEssage 监听Msssage广播消息channel的goroutine,一旦有消息就发送给全部的在线user
func(this *Server)ListenMEssage(){
	for  {
		msg := <-this.Message
		//将消息发送给全部的在线user
		this.mapLock.Lock()
		for _,cli := range this.OnlineMap{
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

// BroadCast 广播消息的方法
func(this *Server)BroadCast(user *User,msg string){
	sendMsg := "["+user.Addr+"]"+user.Name+":"+msg

	this.Message <- sendMsg
}

func (this *Server)Handler(Conn net.Conn){
	//...当前连接的业务
	//fmt.Println("连接建立成功")
	user := NewUser(Conn,this)

	user.Online()
	//监听客户端是否活跃
	isLive := make(chan  bool)
	//接收客户端发送的消息
	go func() {
		buf := make([]byte,4096)
		for{
			n,err := Conn.Read(buf)
			if n == 0 {
				this.BroadCast(user,"下线")
				return
			}
			if err!= nil && err != io.EOF{
				fmt.Println("Conn Read err:",err)
				return
			}

			//提取用户信息
			msg := string(buf[:n-1])

			//用户针对msg进行消息处理
			user.DoMessage(msg)

			//
			isLive <- true
		}
	}()

	//当前handler阻塞
	for{

		select{
		case <- isLive:
			//当前用户活跃，重置定时器
			//不做处理 只是激活select

		case <- time.After(time.Second*10):
			//已经超时
			//将当前客户端强制关闭
			user.SendMsg("你被剔除")
			close(user.C)

			Conn.Close()
			return
		}
	}

}

// Start 启动服务器
func  (this *Server)Start()  {
	//socket listen
	Listener,err := net.Listen("tcp",fmt.Sprintf("%s:%d",this.Ip,this.Post))
	if err != nil{
		fmt.Println("net.Listen err:",err)
		return
	}
	defer Listener.Close()

	go this.ListenMEssage()



	for  {

		//accept
		conn,err := Listener.Accept()
		if err != nil{
			fmt.Println("listen accept:",err)
			continue
		}

		//do handle
		go this.Handler(conn)
		//close listen socket
	}



}