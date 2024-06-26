package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	Message chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}
func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message

		this.mapLock.Lock()
		for _, user := range this.OnlineMap {
			user.C <- msg
		}

		this.mapLock.Unlock()
	}
}
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	this.Message <- sendMsg
}

func (this *Server) Handler(conn net.Conn) {
	defer conn.Close()

	user := NewUser(conn, this)
	user.Online()

	isLive := make(chan bool)

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read Error:", err)
				return
			}

			msg := string(buf[:n-1])

			user.DoMessage(msg)

			isLive <- true
		}
	}()
	for {
		select {
		case <-isLive:
		case <-time.After(time.Minute * 30):
			user.SendMsg("您已超时下线")
			close(user.C)
			conn.Close()
			return
		}
	}

}

func (this *Server) Run() {
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
	}

	if err == nil {
		defer listen.Close()

		go this.ListenMessage()

		for {
			conn, err := listen.Accept()
			if err != nil {
				fmt.Println("listen.Accept err:", err)
				continue
			}

			go this.Handler(conn)

		}
	}
}
