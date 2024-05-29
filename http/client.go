package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial err:", err)
		return nil
	}
	client.conn = conn

	return client
}

func (client *Client) menu() bool {
	var flag int

	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	_, err := fmt.Scanln(&flag)
	if err != nil {
		fmt.Println("Scan err:", err)
		return false
	}

	if flag < 0 || flag > 3 {
		fmt.Println("请输入合法参数")
		return false
	}

	client.flag = flag
	return true
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}

		switch client.flag {
		case 1:
			client.PublicChat()
			break
		case 2:
			client.PrivateChat()
			break
		case 3:
			client.UpdateName()
			break
		}
	}
}

var serverIp string
var serverPort int

func init() {
	fmt.Println("client init")
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "server ip")
	flag.IntVar(&serverPort, "port", 8080, "server port")
}
func main() {
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("client is nil")
		return
	}

	go client.DealResponse()

	fmt.Println("create client successful:", client)

	client.Run()
}

func (client *Client) DealResponse() {
	io.Copy(os.Stdout, client.conn)
}

func (client *Client) UpdateName() bool {
	fmt.Println("请输入用户名:")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name

	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("client.conn.Write err:", err)
		return false
	}

	return true
}

func (client *Client) SelectUsers() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("client.conn.Write err:", err)
		return
	}
}
func (client *Client) PrivateChat() {
	var remoteName string
	var chatMsg string
	client.SelectUsers()
	fmt.Println("请输入聊天对象用户名，exit退出")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println("请输入消息内容，exit退出")
		fmt.Scanln(&chatMsg)
		for chatMsg != "exit" {
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("client.conn.Write err:", err)
					break
				}
			}
			chatMsg = ""
			fmt.Println("请输入消息内容，exit退出")
			fmt.Scanln(&chatMsg)
		}
		client.SelectUsers()
		fmt.Println("请输入聊天对象用户名，exit退出")
		fmt.Scanln(&remoteName)
	}

}
func (client *Client) PublicChat() {
	var chatMsg string
	fmt.Println("请输入聊天内容, exit退出")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		if len(chatMsg) == 0 {
			continue
		}

		sendMsg := chatMsg + "\n"
		_, err := client.conn.Write([]byte(sendMsg))
		if err != nil {
			fmt.Println("client.conn.Write err:", err)
			break
		}

		chatMsg = ""
		fmt.Println("请输入聊天内容, exit退出")
		fmt.Scanln(&chatMsg)
	}
}
