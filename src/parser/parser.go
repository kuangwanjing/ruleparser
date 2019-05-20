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
			operand := exp.Operand
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
		for _, rule := range p.rules[tag] {
			go p.createExamineFn(rule, field.Type.Name(), val.Field(i), ch)()
			count += 1
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

func (p *RuleParser) createExamineFn(rule state.RuleExpr,
	tn string, value reflect.Value, ch chan bool) func() {
	//op := rule.Operation
	// basic data type:https://thorstenball.com/blog/2016/11/16/putting-eval-in-go/

	if !isBasicDataType(tn) {
		if isBasicOperation(rule.Operation) {
			fn := value.MethodByName("Cmp")
			if fn.IsValid() {
				fmt.Println("has Cmp function")
				return func() {
					// use the value of a rule as the argument of customized comparing function
					in := make([]reflect.Value, 1)
					in[0] = reflect.ValueOf(rule.Value)
					ret := fn.Call(in) // how to deal with errors?
					retInt := ret[0].Interface().(int)
					ch <- GetBasicOperation(rule.Operation)(retInt)
				}
			} else {
				fmt.Println("doesn't have Cmp function")
				return func() {
					ch <- false
				}
			}
		} else {
		}
	}

	return func() {
		ch <- true
	}

}
