//Package controller contains the web server for the controller
package controller

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/AscendTech4H/AscendTechROV/go/debug"
	"github.com/AscendTech4H/AscendTechROV/go/startup"
	"github.com/AscendTech4H/AscendTechROV/go/util"

	"github.com/gorilla/websocket"
)

func init() {
	//HTTP startup
	startup.NewTask(200, func() error {
		http.HandleFunc("/websock", websockhandler)
		http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
		return nil
	})
	//Web server startup last
	startup.NewTask(254, func() error {
		log.Println("Starting web controller. . .")
		go http.ListenAndServe(":8080", nil)
		return nil
	})
	lck = new(sync.RWMutex)
}

//conn is a connection used internally
type conn struct {
	sock *websocket.Conn
}

func (c *conn) read() ([]byte, error) {
	_, msg, err := c.sock.ReadMessage()
	if err != nil {
		return nil, err
	}
	return msg, err
}

func (c *conn) write(msg []byte) error {
	return c.sock.WriteMessage(websocket.BinaryMessage, msg)
}

func loadConn(w http.ResponseWriter, r *http.Request) (*conn, error) {
	upgrader := websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	c := new(conn)
	c.sock = connection
	return c, nil
}

var currentConn *websocket.Conn

//Direction constants
const (
	CCW = iota
	STOP
	CW
)

//Robot contains the controller data
type Robot struct {
	Claw, Laser                       bool
	Forward, Up, Tilt, Turn, ClawTurn int
}

var r Robot
var lck *sync.RWMutex //Lock on robot object

//SendData sends data through websocket
func SendData(data []byte) {
	lck.Lock()
	if currentConn != nil {
		currentConn.WriteMessage(websocket.BinaryMessage, data)
	}
	lck.Unlock()
}

func websockhandler(writer http.ResponseWriter, requ *http.Request) {
	connection, err := loadConn(writer, requ)
	util.UhOh(err)
	debug.VLog("Websocket connected")
	defer func() { //If we crash, don't break the robot
		lck.Unlock()
		debug.VLog("Websocket disconnected")
	}()
	for {
		m, e := connection.read() //Read a message
		lck.Lock()
		util.UhOh(e)
		str := string(m)
		if debug.Verbose {
			log.Printf("Websocket Command: %s", str)
		}
		var err error
		switch str[0] {
		case 'C':
			r.Claw = true
		case 'c':
			r.Claw = false
		case 'L':
			r.Laser = true
		case 'l':
			r.Laser = false
		case 'X':
			r.Turn, err = strconv.Atoi(str[1:])
		case 'Y':
			r.Forward, err = strconv.Atoi(str[1:])
		case 'S':
			r.Up, err = strconv.Atoi(str[1:])
		case '{':
			r.ClawTurn = CCW
		case '^':
			r.ClawTurn = STOP
		case '}':
			r.ClawTurn = CW
		case '!':
			r.Tilt, err = strconv.Atoi(str[1:])
		default:
			log.Printf("Ignored unrecognized WebSocket command %s\n", str)
		}
		debug.VLog(str)
		util.UhOh(err)
		lck.Unlock()
	}
}

//RobotState gets robot control state
func RobotState() (rob Robot) {
	lck.RLock()         //Set a read lock
	defer lck.RUnlock() //Unlock on exit
	rob.Claw = r.Claw
	rob.Forward = r.Forward
	rob.Laser = r.Laser
	rob.Turn = r.Turn
	rob.Up = r.Up
	rob.ClawTurn = r.ClawTurn
	return
}
