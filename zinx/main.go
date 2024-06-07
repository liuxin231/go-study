package main

import (
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	s := znet.NewServer()

	s.AddRouter(1, &PingRouter{})

	s.Start()
	// Block, otherwise the listener's goroutine will exit when the main Go exits (阻塞,否则主Go退出， listenner的go将会退出)
	c := make(chan os.Signal, 1)
	// Listen for specified signals: ctrl+c or kill signal (监听指定信号 ctrl+c kill信号)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	sig := <-c
	zlog.Ins().InfoF("[SERVE] Zinx server , name %s, Serve Interrupt, signal = %v", s.ServerName(), sig)
	//s.Serve()
}
