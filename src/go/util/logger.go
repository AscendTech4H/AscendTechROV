package logger

import (
	"fmt"
	"shutdown"
)

var toLog chan string
type Logger struct {								//Logger type
	name string
}

func logroutine() {
	for l := range toLog {
		fmt.Println(l)
	}
}

var debug bool
func SetDebug(deb bool) {							//Set debug mode (off by default)
	debug=deb
}

func init() {
	SetDebug(false)			//Debug defaults to false
	toLog = make(chan string)	//Make logging channel
	go logroutine()			//Start log goroutine
	shutdown.add(stop)
}

func stop() {
	fmt.Println("Stopping logger. . .")
	close(toLog)
}

func NewLogger(nam string) Logger {
	return Logger{name: nam}
}

func brack(txt string) string {						//Put something in brackets
	return "["+txt+"]"
}

func logFmt(t string, name string, msg string) string {			//Format log output
	return brack(t)+brack(name)+" "+msg
}

func (l Logger) Message(msg string) {					//Put a normal message
	toLog <- logFmt("MESSAGE",l.name,msg)
}

func (l Logger) Debug(msg string) {						//Put a debug message (if debug mode is on)
	if debug {
		toLog <- logFmt("MESSAGE",l.name,msg)
	}
}

func (l Logger) ErrorStr(msg string) {						//Put an error message
	toLog <- logFmt("ERROR",l.name,msg)
}

func (l Logger) Error(err error) {						//Put an error message using the error type
	toLog <- logFmt("ERROR",l.name,err.Error())
}
