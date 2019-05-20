package parser

import (
	"fmt"
	"go/scanner"
	"go/token"
	"reflect"
	"state"
	"time"
)

const tagName = "rule"

type RuleParser struct {
	rules     map[string][]state.RuleExpr
	ruleCount int
	timeout   time.Duration
}

func ParserInit(rules string) (*RuleParser, error) {
	return rulesParser(rules)
}

func rulesParser(rules string) (*RuleParser, error) {
	var exprs = make(map[string][]state.RuleExpr)
	var count int = 0

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

		if reflect.TypeOf(curState).Name() == "StateEnd" {
			operand := exp.GetOperand()
			exprs[operand] = append(exprs[operand], exp)
			count += 1
		}

		curState = newState
	}

	rp := &RuleParser{exprs, count, 500 * time.Millisecond}

	return rp, nil
}

func (p *RuleParser) Examine(context interface{}) bool {
	count := 0
	ch := make(chan bool)

	t := reflect.TypeOf(context)
	val := reflect.ValueOf(context)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(tagName)
		if _, ok := p.rules[tag]; ok {
			go p.createExamineFn(tag, val.Field(i), ch)
			count += 1
		}
		vf := val.Field(i)
		fn := vf.MethodByName("Cmp")
		if fn.IsValid() {
			fmt.Println("has Cmp function")
		} else {
			fmt.Println("doesn't have Cmp function")
		}
	}

	for i := 0; i < count; i++ {
		select {
		case rst := <-ch:
			if !rst {
				return false
			}
		case <-time.After(p.timeout):
			return false
		}
	}

	return true
}

func (p *RuleParser) createExamineFn(tag string, value reflect.Value, ch chan bool) func() {
	return func() {
		ch <- true
	}
}
