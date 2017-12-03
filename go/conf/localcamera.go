package conf

import (
	"errors"

	"github.com/AscendTech4H/bracketconf"
)

//LocalCamera is a camera attached to the local computer
type LocalCamera struct {
	Name    string
	DevFile string
}

var localCameraDirectiveProcessor = bracketconf.NewDirectiveProcessor(
	nameDirective,
	bracketconf.Directive{Name: "devfile", Callback: func(object interface{}, ans ...bracketconf.ASTNode) {
		if len(ans) == 1 {
			object.(*LocalCamera).DevFile = ans[0].Text()
		} else if len(ans) == 0 {
			panic(errors.New("Dev file directive has no arguments"))
		} else {
			panic(bracketconf.ConfErr{Pos: ans[0].Position(), Err: errors.New("Dev file directive has too many arguments")})
		}
	}},
)

var localCameraDirective = bracketconf.Directive{Name: "localcamera", Callback: func(object interface{}, ans ...bracketconf.ASTNode) {
	switch len(ans) {
	case 0:
		panic(errors.New("Local camera directive has no arguments"))
	case 1:
		if !ans[0].IsBracket() {
			panic(bracketconf.ConfErr{Pos: ans[0].Position(), Err: errors.New("Single local cam argument must be a bracket")})
		}
		lc := LocalCamera{}
		ans[0].Evaluate(&lc, localCameraDirectiveProcessor)
		object.(*FullList).add(lc)
		return
	case 2:
		object.(*FullList).add(LocalCamera{
			Name:    ans[0].Text(),
			DevFile: ans[1].Text(),
		})
		return
	}
	for _, v := range ans {
		if v.IsBracket() {
			panic(bracketconf.ConfErr{Pos: ans[0].Position(), Err: errors.New("Invalid syntax: cannot mix bracket and non-bracket syntax for a local camera directive")})
		}
	}
	panic(bracketconf.ConfErr{Pos: ans[0].Position(), Err: errors.New("Invalid non-bracket syntax: must have 2 arguments for non-bracket syntax")})
}}
