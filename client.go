package wproxy

import "github.com/gorilla/websocket"

type ClientHandler interface {
	parser                     // 解析器
	handler                    // 处理器
	converter                  // 转换器
	ReadMessage(interface{})   // 客户端原始数据
	HanderError(error)         // 处理异常
	ConnectionError(err error) // 连接异常
	Connected(*websocket.Conn) // 连接成功
}
