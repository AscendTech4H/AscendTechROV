package controller

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"../startup"
	"../util"

	"github.com/gorilla/websocket"
)

func init() {
	//HTTP startup
	startup.NewTask(200, func() error {
		http.HandleFunc("/co_websock", websockhandler)
		http.Handle("/static", http.FileServer(http.Dir("static")))
		return nil
	})
	//Web server startup last
	startup.NewTask(250, func() error {
		log.Println("Starting web controller. . .")
		go http.ListenAndServe(":8081", nil)
		return nil
	})
	lck = new(sync.RWMutex)
}

var currentConn *websocket.Conn

//SendData sends data through websocket
func SendData(data []byte) {
	if currentConn != nil{
		currentConn.WriteMessage(websocket.BinaryMessage, data)
	}
}

func websockhandler(writer http.ResponseWriter, requ *http.Request) {
	var upgrader = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
	connection, err := upgrader.Upgrade(writer, requ, nil)
	util.UhOh(err)
	currentConn = connection //used when sending data
}