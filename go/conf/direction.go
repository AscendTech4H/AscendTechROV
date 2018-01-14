package conf

import (
	"errors"

	"github.com/AscendTech4H/bracketconf"
)

//Axis constants
const (
	XDirection = iota
	YDirection
	ZDirection
)

//Direction is the control settings for a certain direciton
type Direction struct {
	Axis          int
	Motors        []MotAng
	Stabilization string //TODO: change to *ArduinoAccel
}

var directionAxisDirective = bracketconf.Directive{Name: "axis", Callback: func(object interface{}, ans ...bracketconf.ASTNode) {
	if len(ans) == 1 {
		a := ans[0].Text()
		var i int
		if a == "x" {
			i = XDirection
		} else if a == "y" {
			i = YDirection
		} else if a == "z" {
			i = ZDirection
		} else {
			panic(bracketconf.ConfErr{Pos: ans[0].Position(), Err: errors.New("Axis must be X, Y, or Z")})
		}
		object.(*Direction).Axis = i
	} else if len(ans) == 0 {
		panic(errors.New("Axis directive has no arguments"))
	} else {
		panic(bracketconf.ConfErr{Pos: ans[0].Position(), Err: errors.New("Axis directive has too many arguments")})
	}
}}
var directionStabilizationDirective = bracketconf.Directive{Name: "stabilization", Callback: func(object interface{}, ans ...bracketconf.ASTNode) {
	if len(ans) == 1 {
		object.(*Direction).Stabilization = ans[0].Text()
	} else if len(ans) == 0 {
		panic(errors.New("Stabilization directive has no arguments"))
	} else {
		panic(bracketconf.ConfErr{Pos: ans[0].Position(), Err: errors.New("Stabilization directive has too many arguments")})
	}
}}

var directionDirective = bracketconf.Directive{Name: "direction", Callback: func(object interface{}, ans ...bracketconf.ASTNode) {
	if len(ans) == 1 && ans[0].IsBracket() {
		dir := Direction{}
		ans[0].Evaluate(&dir, bracketconf.NewDirectiveProcessor(
			directionAxisDirective,
			motAngDirective,
			directionStabilizationDirective,
		))
		object.(*FullList).add(dir)
	} else if len(ans) == 2 && !ans[0].IsBracket() && ans[1].IsBracket() {

		dir := Direction{}
		a := ans[0].Text()
		var i int
		if a == "x" {
			i = XDirection
		} else if a == "y" {
			i = YDirection
		} else if a == "z" {
			i = ZDirection
		} else {
			panic(bracketconf.ConfErr{Pos: ans[0].Position(), Err: errors.New("Axis must be X, Y, or Z")})
		}
		dir.Axis = i

		ans[1].Evaluate(&dir, bracketconf.NewDirectiveProcessor(
			motAngDirective,
			directionStabilizationDirective,
		))
		object.(*FullList).add(dir)
	} else {
		panic(errors.New("Direction directive needs 1 bracket argument or an int argument and a bracket argument"))
	}
}}
