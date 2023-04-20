package websocket_proxy

import "github.com/gorilla/websocket"

type ServerHandler interface {
	parser                     // 解析器
	filter                     // 过滤器
	handler                    // 处理器
	converter                  // 转换器
	HanderError(error)         // 处理异常
	ReadMessage(interface{})   // 服务端原始数据
	ConnectionError(err error) // 连接异常
	Connected(*websocket.Conn) // 连接成功
}
