package debug

import (
	"flag"
	"html"
	"log"
	"net/http"
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
}

func (l *logBuf) Write(b []byte) (int, error) {
	l.lck.Lock()
	defer l.lck.Unlock()
	s := html.EscapeString(string(b))
	for i, w := range l.outs {
		_, err := w.Write([]byte(s))
		w.(http.Flusher).Flush()
		if err != nil {
			l.outs[i] = nil
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
	return len(b), nil
}
func (l *logBuf) Add(w http.ResponseWriter) {
	l.lck.Lock()
	defer l.lck.Unlock()
	l.outs = append(l.outs, w)
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
		http.HandleFunc("/debug", func(w http.ResponseWriter, r *http.Request) {
			lbuf.Add(w)
		})
		return nil
	})
}
