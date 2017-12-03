package conf

import (
	"errors"
	"net/url"

	"github.com/AscendTech4H/bracketconf"
)

//RemoteCamera is a remotely connected camera
type RemoteCamera struct {
	Name    string
	Address *url.URL
}

var remoteCameraDirectiveProcessor = bracketconf.NewDirectiveProcessor(
	nameDirective,
	bracketconf.Directive{Name: "address", Callback: func(object interface{}, ans ...bracketconf.ASTNode) {
		if len(ans) == 1 {
			u, err := url.Parse(ans[0].Text())
			if err != nil {
				panic(bracketconf.ConfErr{Pos: ans[0].Position(), Err: errors.New("URL given is not valid")})
			}
			object.(*RemoteCamera).Address = u
		} else if len(ans) == 0 {
			panic(errors.New("Address directive has no arguments"))
		} else {
			panic(bracketconf.ConfErr{Pos: ans[0].Position(), Err: errors.New("Address directive has too many arguments")})
		}
	}},
)

var remoteCameraDirective = bracketconf.Directive{Name: "remotecamera", Callback: func(object interface{}, ans ...bracketconf.ASTNode) {
	switch len(ans) {
	case 0:
		panic(errors.New("Remote camera directive has no arguments"))
	case 1:
		if !ans[0].IsBracket() {
			panic(bracketconf.ConfErr{Pos: ans[0].Position(), Err: errors.New("Single remote cam argument must be a bracket")})
		}
		lc := RemoteCamera{}
		ans[0].Evaluate(&lc, remoteCameraDirectiveProcessor)
		object.(*FullList).add(lc)
		return
	case 2:
		u, err := url.Parse(ans[1].Text()) //TODO: does this even work
		if err != nil {
			panic(bracketconf.ConfErr{Pos: ans[1].Position(), Err: errors.New("URL given is not valid")})
		}
		object.(*FullList).add(RemoteCamera{
			Name:    ans[0].Text(),
			Address: u,
		})
		return
	}
	for _, v := range ans {
		if v.IsBracket() {
			panic(bracketconf.ConfErr{Pos: ans[0].Position(), Err: errors.New("Invalid syntax: cannot mix bracket and non-bracket syntax for a remote camera directive")})
		}
	}
	panic(bracketconf.ConfErr{Pos: ans[0].Position(), Err: errors.New("Invalid non-bracket syntax: must have 2 arguments for non-bracket syntax")})
}}
