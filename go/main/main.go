package main

import (
	_ ".."
	"../startup"
)

func main() {
	startup.Start()
	c := make(chan struct{})
	<-c
}
