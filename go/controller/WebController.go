package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
)

func main() {
	http.HandleFunc("/websock", websockhandler)
	http.Handle("/robotcontrol.html", http.FileServer(http.Dir("D:\\gopath\\Site")))
	http.Handle("/virtualjoystick.js", http.FileServer(http.Dir("D:\\gopath\\Site")))
	http.ListenAndServe(":8080",nil)
}
func websockhandler(writer http.ResponseWriter, requ *http.Request) {
	fmt.Print("hai")
	var upgrader = websocket.Upgrader{ReadBufferSize:  1024,WriteBufferSize: 1024,}
	connection, error := upgrader.Upgrade(writer,requ,nil)
	if error != nil {return}
	for {
		t, m, e := connection.ReadMessage()
		if e != nil {return}
		e = connection.WriteMessage(t,m)
		if e != nil {return}
		fmt.Println(m) //use m to control the robot
	}
}