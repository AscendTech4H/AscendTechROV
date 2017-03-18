package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
	"strconv"
)

type robot struct {
	claw, agar, laser bool
	movex,movey,turn int
}
var r robot
func main() {
	http.HandleFunc("/websock", websockhandler)
	http.Handle("/robotcontrol.html", http.FileServer(http.Dir("D:\\gopath\\Site")))
	http.Handle("/virtualjoystick.js", http.FileServer(http.Dir("D:\\gopath\\Site")))
	http.ListenAndServe(":8080",nil)
}
func websockhandler(writer http.ResponseWriter, requ *http.Request) {
	var upgrader = websocket.Upgrader{ReadBufferSize:  1024,WriteBufferSize: 1024,}
	connection, error := upgrader.Upgrade(writer,requ,nil)
	if error != nil {return}
	for {
		t, m, e := connection.ReadMessage()
		if e != nil {return}
		e = connection.WriteMessage(t,m)
		if e != nil {return}
		if (m[0]==[]byte("C")[0]){
			r.claw = true
		} else if (m[0]==[]byte("c")[0]){
			r.claw = false
		} else if (m[0]==[]byte("A")[0]){
			r.agar = true
		} else if (m[0]==[]byte("a")[0]){
			r.agar = false
		} else if (m[0]==[]byte("L")[0]){
			r.laser = true
		} else if (m[0]==[]byte("l")[0]){
			r.laser = false
		} else if (m[0]==[]byte("X")[0]){
			r.movex,_ = strconv.Atoi(string(m[1:]))
		} else if (m[0]==[]byte("Y")[0]){
			r.movey,_ = strconv.Atoi(string(m[1:]))
		} else if (m[0]==[]byte("S")[0]){
			r.turn,_ = strconv.Atoi(string(m[1:]))
		}
		fmt.Println("CLAW:")
		fmt.Println(r.claw)
		fmt.Println("AGAR:")
		fmt.Println(r.agar)
		fmt.Println("LASER:")
		fmt.Println(r.laser)
		fmt.Println("X:")
		fmt.Println(r.movex)
		fmt.Println("X:")
		fmt.Println(r.movey)
		fmt.Println("Turn:")
		fmt.Println(r.turn)
		//now send these to the robot!
	}
}
