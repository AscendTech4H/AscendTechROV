package fake

import (
	"../../controller"
	"log"
)
type fakecontroller struct{
	name string
	buttons map[string]controller.Button
}
type fakejoystick struct{
	name string
	statex uint8
	statey uint8
}
type fakebutton struct{
	name string
	state uint8
}

func (f *fakebutton) State() uint8 {
	return f.state
}
func (f *fakejoystick) State() (uint8, uint8) {
	return f.statex, f.statey
}

func (f *fakebutton) Set(s uint8){
	f.state=s
	log.Printf("Set button %s to %d.\n",f.name,s)
}
func (f *fakejoystick) Set(x uint8, y uint8){
	f.statex=x
	f.statey=y
	log.Printf("Set joystick %s to x: %d y: %d.\n",f.name,x,y)
}

func NewFakeButton(name string, istate uint8) controller.Button {
	m := new(fakebutton)
	m.name=name
	m.state=istate
	return m
}
func NewFakeJoystick(name string, istatex uint8,istatey uint8) controller.Joystick {
	m := new(fakejoystick)
	m.name=name
	m.statex=istatex
	m.statey=istatey
	return m
}
func NewFakeController(name string, buttons map[string]controller.Button) controller.Controller {
	m := new(fakecontroller)
	m.name=name
	m.buttons=buttons
	return m
}

func (f *fakecontroller) GetButtons() map[string]controller.Button{
	return f.buttons
}