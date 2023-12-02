package main

import (
	"net"
	"strings"
)

type  User struct {
	Name string
	Addr string
	C  chan  string
	conn net.Conn

	server *Server
}

// NewUser creat user API
func NewUser(conn net.Conn,server *Server)*User{
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:make(chan string),
		conn: conn,
		server: server,
	}
	go user.ListenMessage()

	return user
}

// Online 用户上线
func (this *User) Online(){

	//用户上线，将用户加入onlinemap中
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	//广播用户上线
	this.server.BroadCast(this,"已上线")


}

// Offline 用户下线
func(this *User) Offline(){
	//用户下线，将用户删除onlinemap中
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap,this.Name)
	this.server.mapLock.Unlock()

	//广播用户上线
	this.server.BroadCast(this,"下线")

}
//对当前客户端发送消息
func (this *User)SendMsg(msg string)  {
	this.conn.Write([]byte(msg))
}

// DoMessage 用户处理消息业务
func(this *User)DoMessage(msg string){
	if msg == "who"{
		//查询当前在线用户
		this.server.mapLock.Lock()
		for _,user := range this.server.OnlineMap{
			onlineMsg := "["+user.Addr+"]"+user.Name+":"+"在线...、\n"
			this.SendMsg(onlineMsg)
		}
		this.server.mapLock.Unlock()
	} else if len(msg)>4 && msg[:3] == "to" {

		//消息格式： to|kidoom|消息内容
		remoteName := strings.Split(msg,"|")[1]
		if remoteName == ""{
			this.SendMsg("消息格式错误")
			return
		}

		remoteUSer , ok := this.server.OnlineMap[remoteName]
		if !ok{
			this.SendMsg("没有该用户")
		}

		content := strings.Split(msg,"|")[2]
		if content == ""{
			this.SendMsg("无消息，请重新发送")
			return
		}
		remoteUSer.SendMsg(this.Name+":"+content)
	}else{
		this.server.BroadCast(this,msg)
	}

}

// ListenMessage listen user channel 一旦有消息就发送给客户端
func (this *User)ListenMessage()  {
	for{
		msg :=<-this.C

		this.conn.Write([]byte(msg+"\n"))
	}

}
