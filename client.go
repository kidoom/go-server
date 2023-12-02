package main

import (
	"flag"
	"fmt"
	"net"
)
var serverIp string
var serverPort int
type Client struct {
	ServerIp string
	ServerPort int
	Name string
	Conn net.Conn
	flag int
}

func Newclient(serverIp string,serverPort int)*Client{
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
	}
	conn,err:=net.Dial("tcp",fmt.Sprintf("%s:%d",serverIp,serverPort))
	if err!=nil{
		fmt.Println("net.Dial error",err)
		return nil
	}
	client.Conn = conn

	return client
}

func (client *Client)menu()bool{
	var flag int
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.退出")
	fmt.Scanln(&flag)
	if flag >= 0 && flag <=3{
		client.flag = flag
		return true
	}else{
		fmt.Println(">>>>>请输入合法范围内的数字<<<<<")
		return false
	}
}
func (Client *Client)Selectuser(){
	sendMsg := "who\n"
	_,err := Client.Conn.Write([]byte(sendMsg))
	if err != nil{
		fmt.Println("conn Write err :",err)
		return
	}
}

func (client *Client)PrivateChat()  {
	var remoteName string
	var chatMsg string
	client.Selectuser()
	fmt.Println(">>>>>请输入聊天对象【用户名】，exit退出")
	fmt.Scanln(&remoteName)

	for remoteName != "exit"{
		if len(chatMsg) != 0{
			sendMSg := "to|"+remoteName+"|"+chatMsg+"\n\n"
			_,err := client.Conn.Write([]byte(sendMSg))
			if err != nil{
				fmt.Println("conn write err:",err)
				break
			}
		}
		chatMsg=""
		fmt.Println(">>>>>请输入聊天内容，exit退出")
		fmt.Scanln(&remoteName)
	}
	client.Selectuser()
	fmt.Println(">>>>>请输入聊天对象【用户名】，exit退出")
	fmt.Scanln(&remoteName)

}

func(client *Client)PublicChat(){
	var chatMsg string
	fmt.Println(">>>>>请输入聊天内容，exit退出")
	fmt.Println(&chatMsg)

	for chatMsg != "exit"{
		if len(chatMsg) != 0{
			sendMsg := chatMsg+"\n"
			_,err := client.Conn.Write([]byte(sendMsg))
			if err != nil{
				fmt.Println("conn write err:",err)
				break
			}
		}
		chatMsg = ""
		fmt.Println(">>>>>请输入聊天内容，exit退出")
		fmt.Println(&chatMsg)
	}

}

//./client -ip 127.0.0.1 -port 8888
func init(){
	flag.StringVar(&serverIp,"ip","127.0.0.1","设置服务器ip地址")
	flag.IntVar(&serverPort,"port",8888,"设置服务器端口")
}

func (client *Client) Run(){
	for client.flag != 0{
		for client.menu() != true{

		}

		switch  client.flag {
		case 1:
			//公聊
			fmt.Println("公聊模式选择")
			client.PublicChat()
			break
		case 2:
			fmt.Println("私聊模式选择")
			client.PrivateChat()
			break
		}

	}
}

func main(){
	flag.Parse()
	client :=Newclient(serverIp,serverPort)
	if client == nil{
		fmt.Println(">>>>> 连接服务器失败<<<<<<")
		return
	}
	fmt.Println(">>>>> 连接服务器成功<<<<<<")
	select {}
}
