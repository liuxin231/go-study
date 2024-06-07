package main

import (
	"fmt"
	"sync"
	"time"
)

const rateLimit = 5

type msgChan struct {
	Id       int64
	Text     string
	DoneChan chan int64
}

var ch chan msgChan

func SendToChan(msg msgChan) {
	ch <- msg
}

func GetFromChan() chan msgChan {
	return ch
}

func producer(i int64) *msgChan {
	var msg msgChan
	msg.Id = i
	msg.Text = fmt.Sprintf("消息文本%v", i)
	msg.DoneChan = make(chan int64)
	SendToChan(msg)
	fmt.Println("producer: ", i)
	return &msg
}

func consumer() {
	for {
		select {
		case msg, ok := <-GetFromChan():
			if ok {
				fmt.Printf("consumer %v processing ...., time: %v\n", msg.Id, time.Now().Format("2006-01-02 15:04:05"))
				time.Sleep(time.Millisecond * (1000 / rateLimit))
				msg.DoneChan <- msg.Id
			}
		}
	}
}

func init() {
	ch = make(chan msgChan, 10)
}
func main() {
	go consumer()

	var wg sync.WaitGroup
	wg.Add(20)

	for i := 0; i < 20; i++ {
		msg := producer(int64(i))
		go func(m *msgChan, w *sync.WaitGroup) {
			defer w.Done()
			if v, ok := <-m.DoneChan; ok {
				fmt.Println("receive done: ", v)
				close(m.DoneChan)
			}
		}(msg, &wg)
	}

	wg.Wait()
	close(ch)
}
