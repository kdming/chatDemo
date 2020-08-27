package main

import (
	"chat/ws"
	"log"
	"net/http"
)

func main() {
	// ws service
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.UpgradeWs(w, r)
	})
	// 打印在线人数
	go ws.ShowOnlineNum()
	// 启动监听服务
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
