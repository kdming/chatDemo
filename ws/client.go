package ws

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// 客户端
type Client struct {
	Key  string
	Conn *websocket.Conn
	Send chan []byte
}

// 消息体
type Message struct {
	Client *Client // 对弈客户端
	Date   string  // 发信日期
	Msg    string  // 消息内容
	Type   string  // 1 普通消息 2 用户上线/离线信息 （可扩展图片信息...)
}

// 消息体工厂函数
func newMsg(c *Client, msg string) *Message {
	m := &Message{}
	m.Type = "1"
	m.Date = time.Now().Format("2006-01-02 15:04")
	m.Client = c
	m.Msg = msg
	return m
}

// 读取信息
func (c *Client) readMsg() {

	defer func() {
		remove(c)
		c.Conn.Close()
	}()

	// 设置读取配置
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {

		// 读取消息
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// 消息内容
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		// 广播信息
		msg := newMsg(c, string(message))
		broadcast(msg)

	}

}

// 发送信息
func (c *Client) sendMsg() {

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {

		select {
		case msg, ok := <-c.Send:

			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// 通知客户端已关闭
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// 创建一个文本类型writer
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(msg)

			// 将排队的信息，字节按顺序发送出去
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			// 发送心跳
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}

	}

}

// 移除无效客户端
func remove(c *Client) {
	delete(Clients, c.Key)
	msg := newMsg(c, c.Key+"下线了！")
	msg.Type = "2"
	broadcast(msg)
}

// 广播消息
func broadcast(msg *Message) {

	for _, client := range Clients {
		if msg.Client.Key != client.Key {
			client.Send <- makeMsg(msg)
		}
	}

}

// 组织消息体
func makeMsg(msg *Message) []byte {

	jStr := `{"user":"%v","msg":"%v", "date":"%v", "type":"%v", "onlineNum": %v}`
	jStr = fmt.Sprintf(jStr, msg.Client.Key, msg.Msg, msg.Date, msg.Type, len(Clients))
	return []byte(jStr)

}
