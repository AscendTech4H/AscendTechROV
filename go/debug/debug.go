//Package debug is used for debugging the robot (via HTML5 SSE or log output)
package debug

import (
	"flag"
	"fmt"
	"html"
	"log"
	"net/http"
	"strings"
	"sync"

	"../startup"
)

//Verbose indicates whether verbose debugging is enabled
var Verbose bool

//Log buffer
type logBuf struct {
	buf  []byte
	outs []http.ResponseWriter
	lck  sync.Mutex
	cle  map[http.ResponseWriter]chan struct{}
}

func (l *logBuf) Write(b []byte) (int, error) {
	l.lck.Lock()
	defer l.lck.Unlock()
	s := "data: " + strings.TrimRight(strings.Replace(html.EscapeString(string(b)), "\n", "<br>", -1), "<br>") + "\n\n"
	for i, w := range l.outs {
		_, err := w.Write([]byte(s))
		w.(http.Flusher).Flush()
		if err != nil {
			l.outs[i] = nil
			l.cle[w] <- struct{}{}
		}
	}
	re := []http.ResponseWriter{}
	for _, w := range l.outs {
		if w != nil {
			re = append(re, w)
		}
	}
	l.outs = re
	l.buf = append(l.buf, b...)
	fmt.Print(string(b))
	return len(b), nil
}
func (l *logBuf) Add(w http.ResponseWriter) chan struct{} {
	l.lck.Lock()
	defer l.lck.Unlock()
	l.outs = append(l.outs, w)
	c := make(chan struct{})
	l.cle[w] = c
	return c
}

var lbuf *logBuf

//VLog - log if verbose
func VLog(s string) {
	if Verbose {
		log.Println(s)
	}
}

func init() {
	startup.NewTask(1, func() error { //Set up can flag parsing
		flag.BoolVar(&Verbose, "verbose", false, "enable verbose debugging")
		return nil
	})
	startup.NewTask(1, func() error {
		lbuf = new(logBuf)
		lbuf.buf = []byte{}
		lbuf.outs = []http.ResponseWriter{}
		lbuf.cle = make(map[http.ResponseWriter]chan struct{})
		http.HandleFunc("/debug", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Connection", "keep-alive")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("data: <b>Connected</b>\n\n"))
			w.(http.Flusher).Flush()
			<-lbuf.Add(w)
		})
		log.SetOutput(lbuf)
		return nil
	})
}
