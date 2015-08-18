package main

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// START UPGRADER OMIT
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// END UPGRADER OMIT

// START READER OMIT
func reader(ws *websocket.Conn) {
	defer ws.Close()
	for {
		if _, _, err := ws.ReadMessage(); err != nil { // HL
			break
		}
	}
}

// END READER OMIT

// START WRITER OMIT
func writer(ws *websocket.Conn) {
	msgTicker := time.NewTicker(5 * time.Second)
	defer func() {
		msgTicker.Stop()
		ws.Close()
	}()

	for {
		select {
		case <-msgTicker.C:
			if err := ws.WriteMessage(websocket.TextMessage, []byte("hello")); err != nil { // HL
				return
			}
		}
	}
}

// END WRITER OMIT

// START SERVE OMIT
func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, _ := upgrader.Upgrade(w, r, nil) // HL
	go writer(ws)
	reader(ws)
}

func main() {
	http.HandleFunc("/ws", serveWs) // HL
	http.ListenAndServe(":8081", nil)
}

// END SERVE OMIT
