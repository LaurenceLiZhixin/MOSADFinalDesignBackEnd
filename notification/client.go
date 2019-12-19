package notification

import (
	"log"
	"net/http"
	"time"
	"github.com/gorilla/websocket"
)

//Client 每个连接上的用户服务
type Client struct {
	HubPtr    *Hub
	Conn      *websocket.Conn
	SendChan  chan []byte
	Useremail string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//NewClientInstance 获取client实例
func NewClientInstance(w http.ResponseWriter, r *http.Request, Thisuseremail string, hub *Hub) (*Client, error) {
	Newconn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("获取client失败")
		log.Println(err)
		return nil, err
	}
	return &Client{
		Conn:      Newconn,
		SendChan:  make(chan []byte),
		Useremail: Thisuseremail,
		HubPtr:    hub,
	}, nil
}

//SendNoti 发送通知函数，将从send chan传入的信息写入conn
func (c *Client) SendNoti() {
	ticker := time.NewTicker(1 *time.Second)
	//断开上下连接
	defer func() {
		c.HubPtr.LogOutChan <- c
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case msg := <-c.SendChan:
			writer, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Println("获取ws的writer失败")
				break
			}
			if _, err := writer.Write(msg); err != nil {
				log.Println("写入ws失败")
				break
			}
			//提高效率，一次发送多条消息
			// n := len(c.sendChan)
			// for i := 0; i < n; i++ {
			// 	if _, err := writer.Write(<-c.sendChan); err != nil {
			// 		log.Panicln("写入ws失败")
			// 		return
			// 	}
			// }
			if err := writer.Close(); err != nil {
				break
			}
		case <-ticker.C:
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
