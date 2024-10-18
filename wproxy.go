package wproxy

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Define ProxyServerOption function type
type ProxyServerOption func(s *ProxyServer)

// Define ProxyServer structure
type ProxyServer struct {
	conn          *websocket.Conn // client connection
	target        *websocket.Conn // server connection
	TargetUrl     string          // proxy target address
	ClientHandler ClientHandler   // client handler
	ServerHandler ServerHandler   // server handler
}

// Parser
type parser interface {
	Parse(interface{}) (interface{}, error)
}

// Filter
type filter interface {
	DoFilter(interface{}) (interface{}, error)
}

// Handler
type handler interface {
	Handler(interface{}) (interface{}, error)
}

// Converter
type converter interface {
	Convert(interface{}) (interface{}, error)
}

// connectIsValid method is used to check if the connection is valid
func (h *ProxyServer) connectIsValid(conn *websocket.Conn) bool {
	if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
		return false
	}
	return true
}

// Proxy method is used to handle WebSocket proxy requests
func (h *ProxyServer) Proxy(w http.ResponseWriter, r *http.Request, cconn, sconn *websocket.Conn, opts ...ProxyServerOption) error {

	var err error

	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // allow WebSocket connections from all sources
		},
	}

	if cconn == nil || !h.connectIsValid(cconn) {
		// Handle request
		cconn, err = upgrader.Upgrade(w, r, nil)
		if err != nil {
			return err
		}
	}
	h.conn = cconn
	// Call the Connected method of ClientHandler
	h.ClientHandler.Connected(cconn)

	for _, opt := range opts {
		opt(h)
	}

	if sconn == nil || !h.connectIsValid(sconn) {
		// Connect to target WebSocket server
		sconn, _, err = websocket.DefaultDialer.Dial(h.TargetUrl, nil)
		if err != nil {
			return err
		}
	}
	h.target = sconn
	// Call the Connected method of ServerHandler
	h.ServerHandler.Connected(sconn)

	var wg sync.WaitGroup
	wg.Add(2)

	// Start two goroutines to handle client and server messages respectively
	go h.handleClientMessages(&wg) // handle client messages
	go h.handleServerMessages(&wg) // handle server messages

	wg.Wait()

	return nil
}

// handleClientMessages handles client messages
func (h *ProxyServer) handleClientMessages(wg *sync.WaitGroup) {

	defer wg.Done()

	for {

		// Read client message
		_, data, err := h.conn.ReadMessage()
		if err != nil {
			// Handle connection error
			h.ClientHandler.ConnectionError(err)
			return
		}
		// Call the ReadMessage method of the client
		h.ClientHandler.ReadMessage(data)

		// Parser
		v, err := h.ClientHandler.Parse(data)
		if err != nil {
			// Handle parsing errors
			h.ClientHandler.HanderError(err)
			continue
		}

		// Handler
		v, messageType, err := h.ClientHandler.Handler(v)
		if err != nil {
			// Handle processing errors
			h.ClientHandler.HanderError(err)
			continue
		}

		// Handle client ping message
		if messageType == PingMessage {
			if err := h.conn.WriteJSON(v); err != nil {
				h.ClientHandler.ConnectionError(err)
				return
			}
			continue
		}

		// Converter
		v, err = h.ClientHandler.Convert(v)
		if err != nil {
			// Handle processing errors
			h.ClientHandler.HanderError(err)
			continue
		}
		// Write to server
		if err := h.target.WriteJSON(v); err != nil {
			// Handle connection error
			if err0 := h.ServerHandler.ConnectionError(err); err0 != nil {
				return
			}
		}
	}

}

// handleServerMessages handles server messages
func (h *ProxyServer) handleServerMessages(wg *sync.WaitGroup) {

	defer wg.Done()

	for {
		// Read server message
		_, data, err := h.target.ReadMessage()
		if err != nil {
			// Handle connection error
			h.ServerHandler.ConnectionError(err)
			return
		}
		// Call the ReadMessage method of the server
		h.ServerHandler.ReadMessage(data)

		// Parser
		v, err := h.ServerHandler.Parse(data)
		if err != nil {
			// Handle parsing errors
			h.ServerHandler.HanderError(err)
			continue
		}
		// Filter
		v, err = h.ServerHandler.DoFilter(v)
		if err != nil {
			// Handle filtering errors
			h.ServerHandler.HanderError(err)
			continue
		}
		if v == nil {
			continue
		}
		// Handler
		v, err = h.ServerHandler.Handler(v)
		if err != nil {
			// Handle processing errors
			h.ServerHandler.HanderError(err)
			continue
		}
		// Converter
		v, err = h.ServerHandler.Convert(v)
		if err != nil {
			// Handle conversion errors
			h.ServerHandler.HanderError(err)
			continue
		}
		// Write to client
		if err := h.conn.WriteJSON(v); err != nil {
			// Handle connection error
			if err0 := h.ClientHandler.ConnectionError(err); err0 != nil {
				return
			}
		}

	}

}
