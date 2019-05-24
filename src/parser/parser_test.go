package parser

import (
	"testing"
)

func TestRulesWithCorrectSyntax(t *testing.T) {
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
		//{"a < 10; b in giergg", "doesn't have `` to quote the value"}, // bugger
	}

	for _, rule := range rules {
		_, err := ParserInit(rule.rule)

		if err == nil {
			t.Errorf("does not detect syntax error `%s`", rule.msg)
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
