package notification

import (
	"errors"
	"log"
)

//Hub 定义通知顶层模块
type Hub struct {
	Clients    map[*Client]bool
	LogInChan  chan *Client
	LogOutChan chan *Client
}

//NewHubInstance 获取顶层模块实例
func NewHubInstance() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		LogInChan:  make(chan *Client),
		LogOutChan: make(chan *Client),
	}
}

//SendMessage 发送通知信息
func (h *Hub) SendMessage(useremail string, msg []byte) error {
	for v := range h.Clients {
		if v.Useremail == useremail {
			v.SendChan <- msg
			return nil
		}
	}
	return errors.New("该用户未在线，发送通知失败")
}

//Run 进行用户的登录登出动态操作
func (h *Hub) Run() {
	for {
		select {
		case tempClient := <-h.LogInChan:
			log.Println("用户注册")
			h.Clients[tempClient] = true
		case tempClient := <-h.LogOutChan:
			log.Println("用户离开")
			delete(h.Clients, tempClient)
		}
	}
}
