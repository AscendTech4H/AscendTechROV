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
		if m[0] == []byte("C")[0] {
			r.Claw = true
		} else if m[0] == []byte("c")[0] {
			r.Claw = false
		} else if m[0] == []byte("A")[0] {
			r.Agar = true
		} else if m[0] == []byte("a")[0] {
			r.Agar = false
		} else if m[0] == []byte("L")[0] {
			r.Laser = true
		} else if m[0] == []byte("l")[0] {
			r.Laser = false
		} else if m[0] == []byte("X")[0] {
			r.Turn, _ = strconv.Atoi(string(m[1:]))
		} else if m[0] == []byte("Y")[0] {
			r.Forward, _ = strconv.Atoi(string(m[1:]))
		} else if m[0] == []byte("S")[0] {
			r.Up, _ = strconv.Atoi(string(m[1:]))
		} else if m[0] == []byte("{")[0] {
			r.ClawTurn = CCW
		} else if m[0] == []byte("^")[0] {
			r.ClawTurn = STOP
		} else if m[0] == []byte("}")[0] {
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
