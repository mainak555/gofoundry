package loggers

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketWriter struct {
	Mutx    sync.Mutex
	Clients map[*websocket.Conn]bool
}

func (wsw *WebSocketWriter) Write(p []byte) (int, error) {
	defer wsw.Mutx.Unlock()
	wsw.Mutx.Lock()
	for con := range wsw.Clients {
		if err := con.WriteMessage(websocket.TextMessage, p); err != nil {
			con.Close()
			delete(wsw.Clients, con)
			fmt.Println("ws-client deleted")
		}
	}
	return len(p), nil
}
