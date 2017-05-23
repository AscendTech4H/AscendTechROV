//Package main is where startup is initiated
package main

import (
	_ "github.com/AscendTech4H/AscendTechROV/go"
	"github.com/AscendTech4H/AscendTechROV/go/startup"
)

func main() {
	startup.Start()
	c := make(chan struct{})
	<-c
}
