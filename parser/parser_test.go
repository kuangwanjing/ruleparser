package parser

import (
	"strconv"
	"strings"
	"testing"
)

func TestRulesWithCorrectSyntax(t *testing.T) {
	rules := []struct {
		rule string
	}{
		{"a < 1"},
		{"a < 10.23"},
		{"bs != false"},
		{"cat == `black,red`"},
		{"a <= 10; b >= 100.3563247; x in `hello,world`"},
		{"x==`10,10,5`;b!=true;t>=-3056"},
	}

	for _, rule := range rules {
		_, err := ParserInit(rule.rule)

		if err != nil {
			t.Errorf("`%s` should be correct", rule.rule)
		}
	}
}

func TestRulesWithIncorrectSyntax(t *testing.T) {
	rules := []struct {
		rule string
		msg  string
	}{
		{"", "empty rule"},
		{";", "contains only semicolumn"},
		{"a <", "doesn't contain value"},
		{"a 10", "doesn't contain operation"},
		{"a < 10, b > 100", "doesn't have semicolumn as a separator"},
		{"a < 10; b in giergg", "doesn't have `` to quote the value"},
	}

	for _, rule := range rules {
		_, err := ParserInit(rule.rule)

		if err == nil {
			t.Errorf("does not detect syntax error `%s`", rule.msg)
		}
	}
}

func TestBasicOperationsForBool(t *testing.T) {
	type TestContext struct {
		T bool `rule:"t"`
	}

	tables := []struct {
		context TestContext
		rules   string
	}{
		{TestContext{true}, "t==true"},
		{TestContext{true}, "t!=false"},
		{TestContext{false}, "t==false"},
		{TestContext{false}, "t!=true"},
		{TestContext{true}, "t==1"},
		{TestContext{true}, "t!=0"},
		{TestContext{false}, "t==0"},
		{TestContext{false}, "t!=1"},
	}

	for _, table := range tables {

		p, err := ParserInit(table.rules)

		if err != nil {
			t.Error("error happens when initializing the parser")
			continue
		}

		rst, err := p.Examine(table.context)

		if err != nil {
			t.Error("error happens when running the parser")
			continue
		}

		if !rst {
			t.Errorf("result of rule `%s` should be true, but false is returned", table.rules)
		}
	}
}

func TestBasicOperationsForString(t *testing.T) {
	type TestContext struct {
		Category string `rule:"category"`
	}

	tables := []struct {
		context TestContext
		rules   string
	}{
		{TestContext{"regular"}, "category == `regular`"},
		{TestContext{"nonregular"}, "category != `regular`"},
		{TestContext{"def"}, "category > `abc`"},
		{TestContext{"abc"}, "category < `defg`"},
		{TestContext{"def"}, "category <=`def`"},
		{TestContext{"def"}, "category <=`defss`"},
		{TestContext{"def"}, "category >=`def`"},
	}

	for _, table := range tables {

		p, err := ParserInit(table.rules)

		if err != nil {
			t.Error("error happens when initializing the parser")
			continue
		}

		rst, err := p.Examine(table.context)

		if err != nil {
			t.Error("error happens when running the parser")
			continue
		}

		if !rst {
			t.Errorf("result of rule `%s` should be true, but false is returned", table.rules)
		}
	}
}

func TestBasicOperationsForInt(t *testing.T) {
	type TestContext struct {
		Age int `rule:"age"`
	}

	tables := []struct {
		context TestContext
		rules   string
	}{
		{TestContext{20}, "age == 20"},
		{TestContext{20}, "age != 22"},
		{TestContext{20}, "age < 22"},
		{TestContext{20}, "age > 15"},
		{TestContext{20}, "age <= 25"},
		{TestContext{20}, "age >= 15"},
		{TestContext{-10}, "age >= -20"},
	}

	for _, table := range tables {

		p, err := ParserInit(table.rules)

		if err != nil {
			t.Error("error happens when initializing the parser")
			continue
		}

		rst, err := p.Examine(table.context)

		if err != nil {
			t.Error("error happens when running the parser")
			continue
		}

		if !rst {
			t.Errorf("result of rule `%s` should be true, but false is returned", table.rules)
		}
	}
}

func TestBasicOperationsForOtherTypes(t *testing.T) {
	type TestContext struct {
		X int8    `rule:"x"`
		Y int64   `rule:"y"`
		Z float32 `rule:"z"`
		W float64 `rule:"w"`
	}

	tables := []struct {
		context TestContext
		rules   string
	}{
		{TestContext{-10, 1559147002, 15.323, 1320.44282}, "x<=0;y>1559137002;z>=15.00;w<2045.348302423"},
		{TestContext{-10, 1559147002, 5.323, 1320.44282}, "x<=0;y>1559137002;z>=15.00;w<2045.348302423"},
		{TestContext{-10, 1559147002, 5.323, 1320.44282}, "x<=0;y>1559137002;z>=-15.00;w<2045.348302423"},
	}

	for i, table := range tables {

		p, err := ParserInit(table.rules)

		if err != nil {
			t.Error("error happens when initializing the parser")
			continue
		}

		rst, err := p.Examine(table.context)

		if err != nil {
			t.Error("error happens when running the parser")
			continue
		}

		if i%2 == 0 && !rst {
			t.Errorf("result of rule `%s` should be true, but false is returned", table.rules)
		}
		if i%2 != 0 && rst {
			t.Errorf("result of rule `%s` should be false, but true is returned", table.rules)
		}
	}
}

func TestNonBasicOperationsForInt(t *testing.T) {
	type TestContext struct {
		Age int `rule:"age"`
	}

	p, err := ParserInit("age between `10, 26`")

	if err != nil {
		t.Error("error happens when initializing the parser")
		return
	}

	_, err = p.Examine(TestContext{20})

	if err == nil {
		t.Error("error should happen when operating non-basic operation for integer field")
	}
}

type TestContext1 struct {
	T TypeT `rule:"t"`
}

type TypeT struct {
	t int
}

func (t TypeT) Cmp(val string) (int, error) {
	i, err := strconv.Atoi(val)

	if err != nil {
		return -1, err
	}

	return t.t - i, nil
}

type TestContext2 struct {
	T *TypeT2 `rule:"t"`
}

type TypeT2 struct {
	t int
}

func (t TypeT2) Cmp(val string) (int, error) {
	i, err := strconv.Atoi(val)

	if err != nil {
		return -1, err
	}

	return t.t - i, nil
}

type TestContext3 struct {
	T *TypeT3 `rule:"t"`
}

type TypeT3 struct {
	t int
}

func (t *TypeT3) Cmp(val string) (int, error) {
	i, err := strconv.Atoi(val)

	if err != nil {
		return -1, err
	}

	return t.t - i, nil
}

type TestContext4 struct {
	T TypeT4 `rule:"t"`
}

type TypeT4 struct {
	t int
}

func TestBasicOperationsForStruct(t *testing.T) {

	p, err := ParserInit("t < 30")

	if err != nil {
		t.Error("error happens when initializing the parser")
		return
	}

	rst, err := p.Examine(TestContext1{TypeT{20}})

	if err != nil {
		t.Error("error happens when operating basic operation for non-basic field")
	}

	if !rst {
		t.Error("should return true")
	}

	rst, err = p.Examine(TestContext2{&TypeT2{20}})

	if err != nil {
		t.Error("error happens when operating basic operation for non-basic field")
	}

	if !rst {
		t.Error("should return true")
	}

	rst, err = p.Examine(TestContext3{&TypeT3{20}})

	if err != nil {
		t.Error("error happens when operating basic operation for non-basic field")
	}

	if !rst {
		t.Error("should return true")
	}

	rst, err = p.Examine(TestContext4{TypeT4{20}})

	if err == nil {
		t.Error("error should happens when operating basic operation for non-basic field without Cmp method")
	}

}

type TestContext5 struct {
	City City `rule:"city"`
}

type City struct {
	val string
}

func (city City) In(val string) (int, error) {
	cities := strings.Split(val, ",")

	for _, c := range cities {
		if city.val == c {
			return 0, nil
		}
	}

	return -1, nil
}

type TestContext6 struct {
	City City2 `rule:"city"`
}

type City2 struct {
	val string
}

func TestNonBasicOperationsForStruct(t *testing.T) {

	p, err := ParserInit("city in `NY,LA,BN,ALT`")

	if err != nil {
		t.Error("error happens when initializing the parser")
		return
	}

	rst, err := p.Examine(TestContext5{City{"NY"}})

	if err != nil {
		t.Error("error happens when operating non-basic operation for non-basic field")
	}

	if !rst {
		t.Error("should return true when operating non-basic operation for non-basic field, but false is returned")
	}

	rst, err = p.Examine(TestContext6{City2{"NY"}})

	if err == nil {
		t.Error("error should happen when operating non-basic operation for non-basic field without defining correponding method")
	}

}

func BenchmarkSum(b *testing.B) {

	type TestContext struct {
		Category string `rule:"category"`
	}

	rules := "category == `regular`"
	context := TestContext{"regular"}

	p, _ := ParserInit(rules)

	for i := 0; i < b.N; i++ {
		p.Examine(context)
	}
}
