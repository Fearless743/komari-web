package terminal

import (
	"sync"

	"github.com/gorilla/websocket"
)

type TerminalSession struct {
	UUID        string
	UserUUID    string
	ClientName  string
	Browser     *websocket.Conn
	Agent       *websocket.Conn
	RequesterIp string
	BytesSent   int64
	BytesRecv   int64
}

var TerminalSessionsMutex = &sync.Mutex{}
var TerminalSessions = make(map[string]*TerminalSession)
