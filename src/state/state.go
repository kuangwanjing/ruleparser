package state

import (
	"errors"
	"fmt"
	"go/token"
)

type RuleExpr struct {
	operand   string
	operation string
	value     string
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
	exp.operand = lit
	return StateOperation{}, nil
}

func (s StateOperation) Run(pos token.Pos, tok token.Token, lit string, exp *RuleExpr) (
	State, error) {
	switch tok {
	case token.IDENT:
		exp.operation = lit
		break
	case token.EQL:
		exp.operation = "=="
		break
	case token.LSS:
		exp.operation = "<"
		break
	case token.GTR:
		exp.operation = ">"
		break
	case token.NEQ:
		exp.operation = "!="
		break
	case token.LEQ:
		exp.operation = "<="
		break
	case token.GEQ:
		exp.operation = ">="
		break
	default:
		return nil, errors.New(fmt.Sprintf("operation is expected at %d", pos))
	}
	return StateValue{}, nil
}

func (s StateValue) Run(pos token.Pos, tok token.Token, lit string, exp *RuleExpr) (
	State, error) {

	if len(lit) <= 2 {
		return nil, errors.New(fmt.Sprintf("operation is expected at %d but empty", pos))
	}

	if lit[0] != '`' || lit[len(lit)-1] != '`' {
		return nil, errors.New(fmt.Sprintf("operation %s is not braced with '`'", lit))
	}
	exp.value = lit[1 : len(lit)-1]
	return StateEnd{}, nil
}

func (s StateEnd) Run(pos token.Pos, tok token.Token, lit string, exp *RuleExpr) (
	State, error) {

	if tok != token.SEMICOLON {
		return nil, errors.New(fmt.Sprintf("`;` is expected at %d", pos))
	}

	return StateOperand{}, nil
}
