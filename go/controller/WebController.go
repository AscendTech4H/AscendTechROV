package controller

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"../startup"

	"github.com/gorilla/websocket"
)

func init() {
	//HTTP startup
	startup.NewTask(200, func() error {
		http.HandleFunc("/websock", websockhandler)
		http.Handle("/static", http.FileServer(http.Dir("static")))
		return nil
	})
	//Web server startup last
	startup.NewTask(250, func() error {
		log.Println("Starting web controller. . .")
		go http.ListenAndServe(":8080", nil)
		return nil
	})
	lck = new(sync.RWMutex)
}

//Direction constants
const (
	CCW = iota
	STOP
	CW
)

//Robot contains the controller data
type Robot struct {
	Claw, Agar, Laser           bool
	Forward, Up, Turn, ClawTurn int
}

var r *Robot
var lck *sync.RWMutex

func websockhandler(writer http.ResponseWriter, requ *http.Request) {
	var upgrader = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
	connection, error := upgrader.Upgrade(writer, requ, nil)
	if error != nil {
		return
	}
	for {
		t, m, e := connection.ReadMessage() //Read a message
		if e != nil {
			return
		}
		e = connection.WriteMessage(t, m) //Do something
		if e != nil {
			return
		}
		lck.Lock()
		str := string(m)
		switch str[0] {
		case 'C':
			r.Claw = true
		case 'c':
			r.Claw = false
		case 'A':
			r.Agar = true
		case 'a':
			r.Agar = false
		case 'L':
			r.Laser = true
		case 'l':
			r.Laser = false
		case 'X':
			r.Turn, _ = strconv.Atoi(string(m[1:]))
		case 'Y':
			r.Forward, _ = strconv.Atoi(string(m[1:]))
		case 'S':
			r.Up, _ = strconv.Atoi(string(m[1:]))
		case '{':
			r.ClawTurn = CCW
		case '^':
			r.ClawTurn = STOP
		case '}':
			r.ClawTurn = CW
		}
		lck.Unlock()
	}
}

//RobotState gets robot control state
func RobotState() (rob Robot) {
	lck.RLock()         //Set a read lock
	defer lck.RUnlock() //Unlock on exit
	rob.Agar = r.Agar
	rob.Claw = r.Claw
	rob.Forward = r.Forward
	rob.Laser = r.Laser
	rob.Turn = r.Turn
	rob.Up = r.Up
	rob.ClawTurn = r.ClawTurn
	return
}
