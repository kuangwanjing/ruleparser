package parser

import (
	"errors"
	"github.com/kuangwanjing/ruleparser/state"
	"go/scanner"
	"go/token"
	"reflect"
	"time"
)

const tagName = "rule"

type RuleParser struct {
	rules     map[string][]state.RuleExpr
	ruleCount int
	timeout   time.Duration
}

type RuleParserChannel struct {
	rst bool
	err error
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

	if len(exprs) == 0 {
		return nil, errors.New("no rules to parse")
	}

	rp := &RuleParser{exprs, count, 500 * time.Millisecond}

	return rp, nil
}

func (p *RuleParser) SetTimeout(t time.Duration) {
	p.timeout = t
}

func (p *RuleParser) Examine(context interface{}) (bool, error) {
	count := 0
	ch := make(chan RuleParserChannel)

	// testing whether context is a pointer. If it is, retrieve the element the pointer points to.
	tmp := reflect.ValueOf(context)
	ck := tmp.Kind()
	if ck == reflect.Ptr || ck == reflect.Interface {
		for ck == reflect.Ptr || ck == reflect.Interface {
			tmp = tmp.Elem()
			ck = tmp.Kind()
		}
		context = tmp.Interface()
		ck = tmp.Kind()
	}

	// since this method handles struct only, it's necessary to determine whether the context is of basic data type.
	// if it is, return an error
	if isBasicDataType(ck.String()) {
		return false, errors.New(ck.String() + " is not accepted")
	}

	t := reflect.TypeOf(context)
	val := reflect.ValueOf(context)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(tagName)
		if tag == "" || tag == "-" {
			continue
		}
		for _, rule := range p.rules[tag] {
			fk := field.Type.Kind()
			fv := val.Field(i)
			go p.createExamineFn(rule, fk, fv, ch)()
			count += 1
		}
	}

	for i := 0; i < count; i++ {
		select {
		case rst := <-ch:
			if !rst.rst || rst.err != nil {
				return rst.rst, rst.err
			}
		case <-time.After(p.timeout):
			return false, errors.New("timeout when parsing")
		}
	}

	return true, nil
}

func (p *RuleParser) createExamineFn(rule state.RuleExpr,
	kind reflect.Kind, value reflect.Value, ch chan RuleParserChannel) func() {

	var k string

	if kind == reflect.Ptr {
		tv := value
		tk := kind
		for tk == reflect.Ptr {
			tv = tv.Elem()
			tk = tv.Kind()
		}
		k = tk.String()
	} else {
		k = kind.String()
	}

	if !isBasicDataType(k) {
		var fnName = ""
		in := make([]reflect.Value, 1)
		in[0] = reflect.ValueOf(rule.Value)

		if isBasicOperation(rule.Operation) {
			fnName = "Cmp"
		} else {
			// convert the first letter into upper case, so that the call is made towards an accessible method
			fnName = ConvertOperationName(rule.Operation)
		}
		fn := value.MethodByName(fnName)
		if !fn.IsValid() {
			fnErr := errors.New(fnName + " function is not found for " + rule.Operand)
			return func() {
				ch <- RuleParserChannel{false, fnErr}
			}
		}
		return func() {
			ret := fn.Call(in)
			retInt, err := p.getReturn(ret)
			if err != nil {
				ch <- RuleParserChannel{false, err}
			}
			if fnName == "Cmp" {
				// basic comparison is built in the package
				ch <- RuleParserChannel{GetBasicOperation(rule.Operation)(retInt), nil}
			} else {
				if retInt != 0 {
					ch <- RuleParserChannel{false, nil}
				} else {
					ch <- RuleParserChannel{true, nil}
				}
			}
		}
	} else {
		if !isBasicOperation(rule.Operation) || isUncomparableDataType(k) {
			fnErr := errors.New(rule.Operation + " is not available for " + rule.Operand)
			return func() {
				ch <- RuleParserChannel{false, fnErr}
			}
		}
		return func() {
			retInt, err := BasicCmp(value.Interface(), rule.Value)
			if err != nil {
				ch <- RuleParserChannel{false, err}
			}
			ch <- RuleParserChannel{GetBasicOperation(rule.Operation)(retInt), nil}
		}
	}

}

func (p *RuleParser) getReturn(ret []reflect.Value) (int, error) {
	if len(ret) != 2 || reflect.TypeOf(ret[0].Interface()).Kind() != reflect.Int {
		return -1, errors.New("function should return an integer and an error object")
	}

	retInt := ret[0].Interface().(int)

	if ret[1].Interface() == nil {
		return retInt, nil
	} else {
		retErr := ret[1].Interface().(error)
		return retInt, retErr
	}
}
