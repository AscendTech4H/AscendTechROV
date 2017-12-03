package conf

import (
	"errors"

	"github.com/AscendTech4H/bracketconf"
)

//MotAng is a motor and an angle
type MotAng struct {
	Name string //motor name
	Ang  float64
}

var motAngDirectiveProcessor = bracketconf.NewDirectiveProcessor(
	nameDirective,
	bracketconf.Directive{Name: "angle", Callback: func(object interface{}, ans ...bracketconf.ASTNode) {
		if len(ans) == 1 {
			object.(*MotAng).Ang = ans[0].Float()
		} else if len(ans) == 0 {
			panic(errors.New("Angle directive has no arguments"))
		} else {
			panic(bracketconf.ConfErr{Pos: ans[0].Position(), Err: errors.New("Angle directive has too many arguments")})
		}
	}},
)

var motAngDirective = bracketconf.Directive{Name: "motorangle", Callback: func(object interface{}, ans ...bracketconf.ASTNode) {
	switch len(ans) {
	case 0:
		panic(errors.New("Motor angle directive has no arguments"))
	case 1:
		if !ans[0].IsBracket() {
			panic(bracketconf.ConfErr{Pos: ans[0].Position(), Err: errors.New("Single motor angle argument must be a bracket")})
		}
		ma := MotAng{}
		ans[0].Evaluate(&ma, motAngDirectiveProcessor)
		object.(*FullList).add(ma)
		return
	case 2:
		object.(*FullList).add(MotAng{
			Name: ans[0].Text(),
			Ang:  ans[1].Float(),
		})
		return
	}
	for _, v := range ans {
		if v.IsBracket() {
			panic(bracketconf.ConfErr{Pos: ans[0].Position(), Err: errors.New("Invalid syntax: cannot mix bracket and non-bracket syntax for a motor angle directive")})
		}
	}
	panic(bracketconf.ConfErr{Pos: ans[0].Position(), Err: errors.New("Invalid non-bracket syntax: must have 2 arguments for non-bracket syntax")})
}}
