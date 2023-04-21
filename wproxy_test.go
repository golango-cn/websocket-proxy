package wproxy

import (
	"log"
	"net/http"
	"testing"

	"github.com/gorilla/websocket"
)

func TestProxy(t *testing.T) {

	// Initialize proxy server
	server := ProxyServer{
		ClientHandler: &clientHandle{},       // Set client handler
		ServerHandler: &serverHandle{},       // Set server handler
		TargetUrl:     "ws://localhost:9090", // Set target URL
	}

	// Proxy websocket service
	http.HandleFunc("/ws/proxy", func(w http.ResponseWriter, r *http.Request) {
		if err := server.Proxy(w, r, func(s *ProxyServer) {
			// Dynamic configuration of target address
			s.TargetUrl = s.TargetUrl + "/?id=123456"
		}); err != nil {
			log.Fatal(err.Error())
		}
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}

}

type clientHandle struct {
}

func (c *clientHandle) Parse(v interface{}) (interface{}, error) {
	// Implementation of parser interface
	// Parse client request data
	return v, nil
}

func (c *clientHandle) Handler(v interface{}) (interface{}, error) {
	// Implementation of handler interface
	return v, nil
}

func (c *clientHandle) Convert(v interface{}) (interface{}, error) {
	// Implementation of converter interface
	return v, nil
}

func (c *clientHandle) ReadMessage(data interface{}) {
	// Original data from client
	// Log record
}

func (c *clientHandle) HanderError(err error) {
	// Implementation of exception handling interface
	// Log record, respond to client with error message, etc.
}

func (c *clientHandle) ConnectionError(err error) {
	// Implementation of connection exception interface
	// Log record, close connection, etc.
}

func (c *clientHandle) Connected(conn *websocket.Conn) {
	// Implementation of connection success interface
	// conn is the client's websocket connection, log record, process logic.
}

type serverHandle struct {
}

func (c *serverHandle) Parse(v interface{}) (interface{}, error) {
	// Implementation of parser interface
	return v, nil
}

func (c *serverHandle) Handler(v interface{}) (interface{}, error) {
	// Implementation of handler interface
	return v, nil
}

func (c *serverHandle) DoFilter(v interface{}) (interface{}, error) {
	// Implementation of filter interface
	// Filter out useless response data from server
	return v, nil
}

func (c *serverHandle) Convert(v interface{}) (interface{}, error) {
	// Implementation of converter interface
	return v, nil
}

func (c *serverHandle) ReadMessage(data interface{}) {
	// Original data from server
	// Log record
}

func (c *serverHandle) HanderError(err error) {
	// Implementation of exception handling interface
	// Log record, etc.
}

func (c *serverHandle) ConnectionError(err error) {
	// Implementation of connection exception interface
	// Log record, close connection, etc.
}

func (c *serverHandle) Connected(conn *websocket.Conn) {
	// Implementation of connection success interface
	// conn is the server's websocket connection, log record, process logic.
}
