package wproxy

import "github.com/gorilla/websocket"

// ClientHandler is an interface that defines the methods for handling WebSocket clients
type ClientHandler interface {
	parser                     // Parser
	handler                    // Handler
	converter                  // Converter
	ReadMessage(interface{})   // Raw data from the client
	HanderError(error)         // Handle exceptions
	ConnectionError(err error) // Connection exceptions
	Connected(*websocket.Conn) // Connection successful
}
