package wproxy

import "github.com/gorilla/websocket"

// ServerHandler is an interface that defines the methods that must be implemented
// by a server handler in order to handle WebSocket connections.
type ServerHandler interface {
	parser                           // Parser
	filter                           // Filter
	handler                          // Handler
	converter                        // Converter
	HanderError(error)               // Handle error
	ReadMessage(interface{})         // Read message from server
	ConnectionError(err error) error // Connection error
	Connected(*websocket.Conn)       // Connection successful
}
