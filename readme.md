简易版的聊天室，只实现了最基本的消息发送与接收，用于练习websocket  
使用 github.com/gorilla/websocket 来实现，可以参考官方demo  
ws 目录为后台代码  
view 目录为前端代码  

启动服务
1. 启动后台服务  
```shell
export GOPROXY=https://goproxy.cn;
go run main.go;
```
2. 启动前端服务
```shell
cnpm i;
npm run serve;
```
