package ws

import (
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var Clients = make(map[string]*Client)

// 协议升级,处理连接
func UpgradeWs(w http.ResponseWriter, r *http.Request) {

	// 获取用户名称
	key := r.URL.Query().Get("key")
	if key == "" {
		w.WriteHeader(401)
		if _, err := w.Write([]byte("身份验证失败")); err != nil {
			log.Println(err)
		}
		return
	}

	// 协议升级
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// 初始化client
	client := &Client{key, conn, make(chan []byte, 256)}
	Clients[key] = client

	go client.readMsg()
	go client.sendMsg()

	// 广播新用户上线信息
	msg := newMsg(client, client.Key+"上线了！")
	msg.Type = "2"
	broadcast(msg)

	ip := RealIP(r)
	log.Println(ip, "连接成功")

}

// 获取连接真实ip
func RealIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ", ")[0]
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	ra, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ra
}

// 打印在线客户端数量
func ShowOnlineNum() {
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ticker.C:
			log.Println("在线客户端数量", len(Clients))

		}
	}
}
