package state

import (
	"errors"
	"fmt"
	"go/token"
)

type RuleExpr struct {
	Operand   string
	Operation string
	Value     string
}

type State interface {
	Run(pos token.Pos, tok token.Token, lit string, exp *RuleExpr) (State, error)
}

type StateOperand struct {
	State
}

type StateOperation struct {
	State
}

type StateValue struct {
	State
}

type StateEnd struct {
	State
}

func (s StateOperand) Run(pos token.Pos, tok token.Token, lit string, exp *RuleExpr) (
	State, error) {
	// examine whether the current token is an identifier
	// if not, return an error description
	if tok != token.IDENT {
		return nil, errors.New(fmt.Sprintf("identifier is expected at %d", pos))
	}
	exp.Operand = lit
	return StateOperation{}, nil
}

func (s StateOperation) Run(pos token.Pos, tok token.Token, lit string, exp *RuleExpr) (
	State, error) {
	switch tok {
	case token.IDENT:
		exp.Operation = lit
		break
	case token.EQL, token.LSS, token.GTR, token.NEQ, token.LEQ, token.GEQ:
		exp.Operation = tok.String()
		break
	default:
		return nil, errors.New(fmt.Sprintf("operation is expected at %d, but %s is found", pos, tok.String()))
	}
	return StateValue{}, nil
}

func (s StateValue) Run(pos token.Pos, tok token.Token, lit string, exp *RuleExpr) (
	State, error) {

	var val string

	if tok == token.STRING {

		if len(lit) <= 2 {
			return nil, errors.New(fmt.Sprintf("operation is expected at %d but empty", pos))
		}

		if lit[0] != '`' || lit[len(lit)-1] != '`' {
			return nil, errors.New(fmt.Sprintf("operation %s is not braced with '`'", lit))
		}

		val = lit[1 : len(lit)-1]
		if val == "" {
			return nil, errors.New(fmt.Sprintf("operation %s is empty", lit))
		}
	} else if tok == token.INT || tok == token.FLOAT {
		val = lit
	} else {
		return nil, errors.New(fmt.Sprintf("%s is not accepted as the value", tok.String()))
	}

	exp.Value = val

	return StateEnd{}, nil
}

func (s StateEnd) Run(pos token.Pos, tok token.Token, lit string, exp *RuleExpr) (
	State, error) {

	if tok != token.SEMICOLON {
		return nil, errors.New(fmt.Sprintf("`;` is expected at %d", pos))
	}

	return StateOperand{}, nil
}
