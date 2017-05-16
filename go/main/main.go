//Package main is where startup is initiated
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
