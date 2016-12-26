package shutdown

type ShutdownHandler func()

var sdh []ShutdownHandler

func init() {
	sdh = make([]ShutdownHandler, 0)
}

func Add(s ShutdownHandler){
	sdh = append(sdh,s)
}

func Shutdown() {
	for _,f := range sdh {
		f()
	}
}
