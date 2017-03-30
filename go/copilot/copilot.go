package copilot

import (
	"net/http"
	"sync"

	"../startup"
	"../util"

	"github.com/gorilla/websocket"
)

var lck *sync.Mutex

func init() {
	//HTTP startup
	startup.NewTask(200, func() error {
		http.HandleFunc("/co_websock", websockhandler)
		return nil
	})
	lck = new(sync.Mutex)
}

var currentConn *websocket.Conn

//SendData sends data through websocket
func SendData(data []byte) {
	lck.Lock()
	defer lck.Unlock()
	if currentConn != nil {
		currentConn.WriteMessage(websocket.BinaryMessage, data)
	}
}

func websockhandler(writer http.ResponseWriter, requ *http.Request) {
	lck.Lock()
	defer lck.Unlock()
	var upgrader = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
	connection, err := upgrader.Upgrade(writer, requ, nil)
	util.UhOh(err)
	currentConn = connection //used when sending data
}
