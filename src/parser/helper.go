package parser

var basicOperations []string = []string{
	"==",
	"!=",
	"<",
	"<=",
	">",
	">=",
}

var basicDataTypes []string = []string{
	"int",
	"string",
}

func isBasicOperation(op string) bool {
	for _, bo := range basicOperations {
		if op == bo {
			return true
		}
	}
	return false
}

func isBasicDataType(tn string) bool {
	for _, bt := range basicDataTypes {
		if tn == bt {
			return true
		}
	}
	return false
}

func GetBasicOperation(op string) func(int) bool {
	switch op {
	case "==":
		return BasicEqual
	case "!=":
		return BasicNotEqual
	case "<":
		return BasicLessThan
	case "<=":
		return BasicLessThanOrEqual
	case ">":
		return BasicGreaterThan
	case ">=":
		return BasicGreaterThanOrEqual
	}
	return nil
}

func BasicEqual(cmp int) bool {
	return cmp == 0
}

func BasicNotEqual(cmp int) bool {
	return cmp != 0
}

func BasicGreaterThan(cmp int) bool {
	return cmp > 0
}

func BasicGreaterThanOrEqual(cmp int) bool {
	return cmp > 0 || cmp == 0
}

func BasicLessThan(cmp int) bool {
	return cmp < 0
}

func BasicLessThanOrEqual(cmp int) bool {
	return cmp < 0 || cmp == 0
}
