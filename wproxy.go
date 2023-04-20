package wproxy

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// 定义 ProxyServerOption 函数类型
type ProxyServerOption func(s *ProxyServer)

// 定义 ProxyServer 结构体
type ProxyServer struct {
	conn          *websocket.Conn // 客户端连接
	target        *websocket.Conn // 服务端连接
	TargetUrl     string          // 代理目标地址
	ClientHandler ClientHandler   // 客户端处理
	ServerHandler ServerHandler   // 服务端处理
}

// 解析器
type parser interface {
	Parse(interface{}) (interface{}, error)
}

// 过滤器
type filter interface {
	DoFilter(interface{}) (interface{}, error)
}

// 处理器
type handler interface {
	Handler(interface{}) (interface{}, error)
}

// 转换器
type converter interface {
	Convert(interface{}) (interface{}, error)
}

// Proxy 方法用于处理 WebSocket 代理请求
func (h *ProxyServer) Proxy(w http.ResponseWriter, r *http.Request, opts ...ProxyServerOption) error {

	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // 允许所有来源的 WebSocket 连接
		},
	}

	// 处理请求
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	h.conn = conn

	// 调用ClientHandler的连接成功方法
	h.ClientHandler.Connected(conn)

	for _, opt := range opts {
		opt(h)
	}

	// 连接到目标 WebSocket 服务器
	target, _, err := websocket.DefaultDialer.Dial(h.TargetUrl, nil)
	if err != nil {
		return err
	}
	h.target = target
	h.ServerHandler.Connected(target)

	var wg sync.WaitGroup
	wg.Add(2)

	// 启动两个 goroutine，分别处理客户端和服务端的消息
	go h.handleClientMessages(&wg) // 处理客户端消息
	go h.handleServerMessages(&wg) // 处理服务端消息

	wg.Wait()

	return nil
}

// handleClientMessages 处理客户端消息
func (h *ProxyServer) handleClientMessages(wg *sync.WaitGroup) {

	defer wg.Done()

	for {
		// 读取客户端消息
		_, data, err := h.conn.ReadMessage()
		if err != nil {
			// 连接错误处理
			h.ClientHandler.ConnectionError(err)
			return
		}
		// 调用客户端消息读取方法
		h.ClientHandler.ReadMessage(data)

		// 解析器
		v, err := h.ClientHandler.Parse(data)
		if err != nil {
			// 处理解析错误
			h.ClientHandler.HanderError(err)
			continue
		}
		// 处理器
		v, err = h.ClientHandler.Handler(v)
		if err != nil {
			// 处理处理错误
			h.ClientHandler.HanderError(err)
			continue
		}
		// 转换器
		v, err = h.ClientHandler.Convert(v)
		if err != nil {
			// 处理处理错误
			h.ClientHandler.HanderError(err)
			continue
		}
		// 写入服务端
		if err := h.target.WriteJSON(v); err != nil {
			// 连接错误处理
			h.ServerHandler.ConnectionError(err)
			return
		}
	}

}

// handleServerMessages 处理服务端消息
func (h *ProxyServer) handleServerMessages(wg *sync.WaitGroup) {

	defer wg.Done()

	for {
		// 读取服务端消息
		_, data, err := h.target.ReadMessage()
		if err != nil {
			// 连接错误处理
			h.ServerHandler.ConnectionError(err)
			return
		}
		// 调用服务端消息读取方法
		h.ServerHandler.ReadMessage(data)

		// 解析器
		v, err := h.ServerHandler.Parse(data)
		if err != nil {
			// 处理解析错误
			h.ServerHandler.HanderError(err)
			continue
		}
		// 过滤器
		v, err = h.ServerHandler.DoFilter(v)
		if err != nil {
			// 处理过滤错误
			h.ServerHandler.HanderError(err)
			continue
		}
		if v == nil {
			continue
		}
		// 处理器
		v, err = h.ServerHandler.Handler(v)
		if err != nil {
			// 处理处理错误
			h.ServerHandler.HanderError(err)
			continue
		}
		// 转换器
		v, err = h.ServerHandler.Convert(v)
		if err != nil {
			// 处理转换错误
			h.ServerHandler.HanderError(err)
			continue
		}
		// 写入客户端
		if err := h.conn.WriteJSON(v); err != nil {
			// 连接错误处理
			h.ClientHandler.ConnectionError(err)
			return
		}
	}

}
