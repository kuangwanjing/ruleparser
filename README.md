# ruleparser

Package ruleparser provide a generic tool to filter data with predefined conditions, which are called rules here.  For example, if a cross-platform application is required to handle the variance of platform, software version, os version, the application provider can construct different variants of configuration with conditions denoting what platform the configuration applies for. Hence, a rule-parser which understands the rules and the current context and determines whether the context matches up with the rules, is needed. This is what ruleparser does. 

## Installment

download the package:

```shell
go get github.com/kuangwanjing/ruleparser
```

## How to use the parser? 4 steps:

### Step 1: Input the rules to examine and get a parser. 

```go
import (
  "github.com/kuangwanjing/ruleparser/parser"
)

rules := "platform == `android`;version < `1.3.2`;field1 > 10; field2 in `val1,val2,val3`"
p, err := parser.ParserInit(rules) // if any error is met during the build of the parser, an error is return. 
```

#### Syntax of the rules

Package ruleparser accepts a bunch of binary conditions separated by semicolon ";". Each rule contain three entities — an operand, a binary operation and matching value separating by space. 

The **operand** is just a legal identity in Go. 

The **operation** is legal when its value falls into one of the 7 categories : 

```Go
== 	// 1. equality
!= 	// 2. inequality
> 	// 3. greater than
>= 	// 4. greater than or equal to
< 	// 5. less than
<= 	// 6. less than or equal to
ide	// 7. or any other legal identity.
```

This means the parser is capable to handle mathematical likewise comparison and provides customized operation for application's needs.  

The value should be designed to match up with the operand. It can be any integer, float number, string or any serialized pattern in string like the above example. 

Values of bool, int, int8, int16, int32, int64,uint, uint32,uint64,string,float32,float64 are legal as value in the rule. **Any string should be enclosed by "`".**

### Step 2: Define the struct for context, point out the struct tags for parsing

In this package, a special struct tag "rule" is used to point out the fields of a struct to be parsed and the struct tags are used to map the field with specific rules. For example:

```go
// define struct for context SoftwareInfo
type SoftwareInfo struct {
  Ver 				Version `rule:"ver"`
  Platform 		string 	`rule:"platform"`
  OtherField1 int 		`rule:"field1"`
  OtherField2 Field2 	`rule:"field2"`
}

type Version struct {
  val string
}

type Field2 struct {
  val string
}

//Initialize the context data
software := SoftwareInfo{Version{"1.3.5"},"android",20,Field2{"val3"}}
```

Basic data type like int and string can be compared if the operation is one of  "==", "!=", ">", ">=", "<", "<=". For example, the field "Platform" and "OtherField1" can be compared with rules platform == \`android\` and field1 > 10 by the parser.

### Step 3: Define customized operation

#### I. Value-Comparison operation with object method. 

For example, the versions of software can be compared by their values. Instead of turning it into integers or string, a correct way to compare them is to comparing each level of the version number. Therefore, we need to define a method "Cmp" to compare versions. If Cmp method is not defined for this data type, an error is returned when parsing.

```go
type Version struct {
  val string
}
func (ver Version) Cmp (pattern string) (int, error) {
  //... 
  // if input pattern is equal to ver
  return 0, nil // return 0
  // if input > ver
  return 1, nil // return a positive integer
  // if input < ver
  return -1, nil // return a negative integer
  // return any error if the pattern is not in correct form
}
```

#### II. Non-Value-Comparison operation with object method. 

When we need to define other operation to check the rules for the application but simple value-comparision can not satisfy the need, we can define a method with similar Name of the rule's operation. For example, we have a rule "field2 in \`val1,val2,val3\`", so we need to define "In" method for data type Field2. (The parser searches for the method with a name of converting the first letter of the operation into upper case so that the searched method is exposed to the parser to invoke.)

```go
type Field2 struct {
  val string
}

func (f Field2) In (pattern string) (int, error) {
  // if the value of f is match with the pattern
  return 0
  // otherwise return a non-zero integer
  return -1
  // for example
  vs := strings.Split(pattern string)
  for _, v := range vs {
    if f.val == v {
      return 0, nil
    }
  }
  return -1, nil
}
```

This is very useful when the operation is about regular expression matching, a remote procedure checking. 

### Step 4: Run the parser

```go
// input the software object 
rst, err := p.Examine(software) // where rst is a bool representing whether the object matches with the rule and err is error.
// or input the address of the object
rst, err := p.Examine(&software)
if err != nil {
  // we got some errors here ... 
}
```

#### Here listing possible errors from the parser

1. The value of the rule doesn't match with the field of the context. For example, if a field is defined as integer but a float number is represent in the rule, the comparison would go wrong. 
2. The comparison method is missing. For example, the missing of Cmp method for field Ver or  In method for field OtherField2 leads to a method missing error.
3. The comparison method is timeout. If the method is dealing with a RPC and the request goes timeout, the parser aborts the whole procedure and returns an error. The default timeout time is 500ms but is changeable via SetTimeout method of the parser which accepts a time.Duration object. 
4. Input a non-object value as the Examine argument. For example, input an integer to the parser.
5. Binding the comparison method with incorrect receiver. For example, if the "Cmp" method is bound to a pointer receiver but the field defined in the struct is in value style, the parser can not detect the method. For example:

```go
type SoftwareInfo struct {
  Ver 				Version `rule:"ver"`
  Platform 		string 	`rule:"platform"`
  OtherField1 int 		`rule:"field1"`
  OtherField2 Field2 	`rule:"field2"`
}

type Version struct {
  val string
}

func (ver *Version) Cmp(pattern string) (int, error) {
  // !! this method would not found by the parser
}
```

Therefore it is safe to define the field and define the method in corresponding receiver (whether it is a value or a pointer receiver)

More example can be found in the example directory. 
