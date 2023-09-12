package wproxy

import "github.com/gorilla/websocket"

type ClientMessageType string

const NormalMessage ClientMessageType = "Normal"
const PingMessage ClientMessageType = "Ping"

// ClientHandler is an interface that defines the methods for handling WebSocket clients
type ClientHandler interface {
	parser                                                       // Parser
	Handler(interface{}) (interface{}, ClientMessageType, error) // Handler
	converter                                                    // Converter
	ReadMessage(interface{})                                     // Raw data from the client
	HanderError(error)                                           // Handle exceptions
	ConnectionError(err error) error                             // Connection exceptions
	Connected(*websocket.Conn)                                   // Connection successful
}
