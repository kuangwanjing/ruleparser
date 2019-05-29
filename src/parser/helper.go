package parser

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

var basicOperations []string = []string{
	"==",
	"!=",
	"<",
	"<=",
	">",
	">=",
}

var basicDataTypes []string = []string{
	"bool",
	"int",
	"int8",
	"int16",
	"int32",
	"int64",
	"uint",
	"uint8",
	"uint16",
	"uint32",
	"uint64",
	"uintptr",
	"string",
	"array",
	"map",
	"slice",
	"ptr",
	"array",
	"float32",
	"float64",
	"interface",
	"chan",
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

func isUncomparableDataType(tn string) bool {
	return tn == "array" || tn == "map" || tn == "slice" || tn == "ptr" ||
		tn == "interface" || tn == "chan"
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
	return func(int) bool {
		return false
	}
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

func ConvertOperationName(op string) string {
	return strings.ToUpper(string(op[0])) + op[1:]
}

// not implemented
func BasicCmp(val interface{}, cmpVal string) (int, error) {

	k := reflect.TypeOf(val).Kind()

	switch k {
	case reflect.String:
		return strings.Compare(val.(string), cmpVal), nil
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(cmpVal)
		if err != nil {
			return -1, err
		}
		if val.(bool) == boolVal {
			return 0, nil
		} else {
			return 1, nil
		}
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int, reflect.Int64:
		int64Val, err := strconv.ParseInt(cmpVal, 10, 64)
		if err != nil {
			return -1, err
		}
		v := reflect.ValueOf(val).Int()
		if v == int64Val {
			return 0, nil
		} else if v < int64Val {
			return -1, nil
		} else {
			return 1, nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uint64Val, err := strconv.ParseUint(cmpVal, 10, 64)
		if err != nil {
			return -1, err
		}
		v := reflect.ValueOf(val).Uint()
		if v == uint64Val {
			return 0, nil
		} else if v < uint64Val {
			return -1, nil
		} else {
			return 1, nil
		}
	case reflect.Float32, reflect.Float64:
		float64Val, err := strconv.ParseFloat(cmpVal, 64)
		if err != nil {
			return -1, err
		}
		v := reflect.ValueOf(val).Float()
		if v == float64Val {
			return 0, nil
		} else if v < float64Val {
			return -1, nil
		} else {
			return 1, nil
		}
	}

	return 0, errors.New("type error for " + cmpVal)
}
