package parser

import (
	"go/scanner"
	"go/token"
	"reflect"
	"state"
	//"fmt"
)

type RuleParser struct {
	rules []state.RuleExpr
}

func ParserInit(rules string) (RuleParser, error) {
	exprs, err := rulesParser(rules)
	if err != nil {
		return RuleParser{}, err
	} else {
		rp := RuleParser{exprs}
		return rp, nil
	}
}

func rulesParser(rules string) ([]state.RuleExpr, error) {
	var exprs []state.RuleExpr

	// Initialize the scanner.
	var s scanner.Scanner
	fset := token.NewFileSet()                        // positions are relative to fset
	file := fset.AddFile("", fset.Base(), len(rules)) // register input "file"
	s.Init(file, []byte(rules), nil /* no error handler */, scanner.ScanComments)

	// Repeated calls to Scan yield the token sequence found in the input.
	var curState state.State = state.StateOperand{}
	var exp state.RuleExpr
	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}

		if reflect.TypeOf(curState).Name() == "StateOperand" {
			exp = state.RuleExpr{}
		}

		newState, err := curState.Run(pos, tok, lit, &exp)

		if err != nil {
			return nil, err
		}

		//fmt.Printf("%s,%s\n", reflect.TypeOf(curState).Name(), reflect.TypeOf(exp))

		if reflect.TypeOf(curState).Name() == "StateEnd" {
			exprs = append(exprs, exp)
		}

		curState = newState
	}

	return exprs, nil
}
