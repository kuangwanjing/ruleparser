## Design Document — Application Rules Parser implemented with Golang

- Background
- Data Abstraction 
- Design and Implementation
- Problem and Tradeoff
- Reference

### Background

Many applications need to take action under pre-defined rules. For example, a subscriber of videos website is able to block out some channels or some topics from their recommendation list. The server of the website is responsible to rule out those videos according to the subscriber's favor. The configuration of an application needs to have variations among different platform, OS version, software version, etc. This means when generating application configuration, one needs to apply some particular rules to choose the correct version of configuration. 

This post is going to discuss a general solution to understand those rules and run the rules automatically so that the wanted data is achieved while enabling extensible rules to affect the result so that new rules can be elegantly added to the application. 

### Data Abstraction

In this context, we will have three types of entities — Data, rules and context. 

- Data: Data is the entity that reach the clients and serve the clients. 
- Rules: Here rules tell whether data is available to the clients. Abstraction of rules is a binary operation in form of "operand operator value". For example, if we want to have a rule describing a subscriber wanting to filter videos from movie channels, then we can have a rule like this: "channel != movie".
- Context: Context tells about data. Take the above example here again. If there is a video provider that provide video containing movies, advertisement, education materials, then we have a context with channels of the above sources. By applying the rules, some part of the context is not acceptable, which means that data corresponding to that context is not acceptable by the clients. 

The following diagram tells us how this data-rules-context model work to serve the favored video to the clients:

![data structure](/Users/kuangwanjing/Documents/code/rule_parser_ds.png)

Note that in different applications, rules and context may lie in different software layers. To be clear, in the video subscribing application, the context can be just the attributions of videos and rules lie in the user's profile. On the other hand, in the example of application configuration management, the context represents the application runtime environment and the rules lie in the description of a variation of the configuration. Therefore, in this design, these three concepts are separated as independent entity to make it more understandable. 

### Design and Implementation

So far, we know that data is just a bunch of objects to be examined, rules are a series of binary operations and context is key-value pairs. Therefore, we need to:

1. Define a straighforward syntax to represent rules. Make sure the syntax is generic enough for genenal purposes, generation simplicity and satifies certain storage efficiency. 
2. Construct a parser to understand the rules and the context. Besides, generate examining codes to examine the rules and context efficiently. 

#### Rule Syntax

As mentioned, the rule syntax is contained with a variable respresenting which part of context to be examined, an operation corresponding to an action on the context, and a value or more precisely a pattern the context should fall into. 

`$val op $pattern` is the syntax for rules. It turns out to be generic enough for our application: operations can be totally reinvented from the scratch and managed by application itself; pattern can be a string containing anything. Thus, a very specific operation is available for representation. For example "the location of the user should be 5 miles within some x,y coordination" can be respresents as `$loc in (x,y,5)`. 

As rules should be persistent somewhere in the application, so storage efficiency is under consideration. The rules can be compacted by either these ways: compacting the $val with short ids and managing the mapping from ids to real variables; compacting the operations in the same way.

#### Rule Parser

The input of the parser is a bunch of rules and the context. First, we need to extract every piece of information from the context and parse every rule to match with the context. Therefore, to run the rule parser, we need two things — context extraction and analysis of lexer and syntax of the rules.

For lexical and syntax analysis, a lexer scanner and state machine to build a rule syntax tree are needed. In fact, as soon as the rule is tokenized and the syntax tree is built, a corresponding matching action can be taken and another rule can be parsed at the same time to improve the performance. 

![](/Users/kuangwanjing/Documents/code/rule_parser.png)

The rule parser output several rule executors which examine whether the particular field of the context falls into the pattern provided by the rule. The syntax of the rule is simple so that lexical scanning, state machine  and parsing are equivalent to parsing a mathematical expression like "c = a + b", which is solvable. The remaining question is how to generate the rule executors. How a rule executor behaves is determined by the rule and the context itself. If the rule is going to examine the version of the software, then at least the executor should have understand the format of the version. Therefore, by the understanding of the fields of the context, executors can be generated: it works like operation overloading in OOP and the behavior of the operation is determined by the type of the caller. Hence, we can bind the context with the examination behaviors and extract the proper behavior for each rule from the corresponding field of context to generate rule executors. That's why we need to extract information about the context. 

#### Implementation

##### Context Extraction

First of all, extract the fields of context so that the binded behaviors can be extracted. This can be achieved by the reflection mechanism of Golang. Although different application contains different context fields, the context is defined as Golang's struct and the struct can be reflected through reflection package so that the type, name and descriptions of fields can be extracted. But how to match a field with a particular rule, given a rule, how do we know that this rule concerns with that field of context? The name of the field is an available choice because it can link the name of operation in the rule with the name of the field in the context. But it is not the best one since compression for the rule can happens. So we need another way for matching. The traditional solution provided by Golang is to use the struct tag. This is the way it does for json or xml serialization or deserialization. Suppose we have the following the data structure for the context:

```go
type SoftwareInfo struct {
  Ver Version `rule:"ver"`,
  ReleaseDate time.Time `rule:"ts"`
}
```

Here we define a special struct tag field called "rule" for the parser's usage(for json usage, that is "json" and "xml" for xml). If this structure needs to also have json operations, then add json struct tag to it separated by comma and this won't conflict with the parsing struct tags. 

```go
type SoftwareInfo struct {
  Ver Version `rule:"ver",json:"version"`,
  ReleaseDate time.Time `rule:"ts",json:"release_date"`
}
```

By iterating the tags, we can have a mapping between the context fields and the rules. For example, we know that rule "ver < 1.2.0" is talking about Ver field of SoftwareInfo. There is another advantage we will have by using struct tag, that is the capability of the parser can be extended by adding different descriptors inside the struct tags. 

##### Lexical Scanning

This process contains two main parts — lexical analysis, syntax analysis. For lexical analysis, Golang provides a useful package — scanner to extract the token of lines of code-like text and provides the literal meaning, token type of a token. For example, we have a string "ver < \`1.2.0\`". Then by scanner package, we extract the tokens from the string.

```go
package main

import (
  "go/scanner"
  "go/token"
)

func main() {
  src := []byte("ver < `1.2.0`")
	var s scanner.Scanner
	fset := token.NewFileSet()                      // positions are relative to fset
	file := fset.AddFile("", fset.Base(), len(src)) // register input "file"
	s.Init(file, src, nil /* no error handler */, scanner.ScanComments)
  // Repeated calls to Scan yield the token sequence found in the input.
	for {
		_, tok, lit := s.Scan() // the first value of the return is the position for the token
		if tok == token.EOF {
			break
		}
    // ...
    // tok can be an identifier, an operation symbol, a string, an integer ...
  }
}

```

For simplicity of the state machine discussed later, the pattern of the rule is quoted by "``"  so that the whole pattern is recognized by the scanner as a single token, otherwise the above "1.2.0" would be recognized as two float numbers. 

##### State Machine

The state machine for this application contains only three states: Identifier —> Operation —> Pattern.  Once we get a token from the scanning, we can put it into the state machine.

##### Parsing

Once we have extracted the identifier, operation and pattern of a rule, parsing is just to combine these three to generate a function operated on corresponding field with the pattern. How the function works depends on the operation and the data type of identifier(field). As the caller of the rule parser, it needs to define the operation handler. Take the SoftwareInfo as the example again. 

```go
type SoftwareInfo struct {
  Ver Version `rule:"ver"`,
  ReleaseDate time.Time `rule:"ts"`
}

// define operation hanlder for version comparison
func (ver Version) Cmp(v Version) {
	// do the version comparison logic here ...
  return 0 // -1  or 1
}
```

Again by reflection, the runtime data type of the interface is understood and can call its own method by method's name. The method's name is determined by the operation.  

```go
var context = SoftwareInfo(/*set the attributes here*/)
t := reflect.TypeOf(context)
val := reflect.ValueOf(&context).Elem()
for i := 0; i < t.NumField(); i++ {
  if (/*getting version field from struct tag*/) {
    vf := val.Field(i).Addr()
    ret := vf.MethodByName("Cmp").Call(/* we need to input the pattern here*/)
    if (ret) { /*examine the return value*/
      // ...
    } else {
      // ...
    }
  }
}
```

Since there are so many [tokens](https://golang.org/src/go/token/token.go?s=422:436#L3) including "<=", ">=", "<", ">", "!=", "==", the parser should have the ability to deal with general comparison operations. Here are some principles of the parser:

1. If the data type of the field is basic type like Integer or String, those comparision can be done without the guide of the caller. 
2. If the operations of a field only contain basic comparisons, the caller only needs to provide "Cmp" guidance. 
3. If the caller needs to define a new operation, define a method with the same name or provide a method mapping. This method should return 0 if the operation succeeds or non-zero otherwise. 
4. It's the caller to deal with the pattern. 

Next is to generate rule executors. We can use the functional programming style to encapsulate a function generator which take the method of the field and the pattern as the input. 

```go
func RuleExecutorGenerate(field interface{}, operation string, pattern string) func() int {
  method := // ... extract method by reflection with operation
  return func() int {
    return method.Call(pattern)
  }
}
```

Now we have rule executor for each rule, we can encapsulate the executors as a comprehensive parser for reuse purpose because sometimes we want the same set of rules to examine a series of data with different context. Therefore, the usage of the parser should be like this:

```go
pser := parser.Parser(context, rules) // context is a prototype, rules are array of strings
for d in data {
  rst, msg := pser.examine(d.context) 
  if rst {
    // the context passes
  } else {
    // msg is the description of the reason the context fails the rules examination. 
  }
}
```

### Reference

[Go Token](https://golang.org/src/go/token/token.go)

[Go Scanner](https://golang.org/pkg/go/scanner/)

[Go Reflection](https://golang.org/pkg/reflect/) and [Laws of Reflection](https://blog.golang.org/laws-of-reflection)

[Lexical Scanner by Rod Pike](https://www.youtube.com/watch?v=HxaD_trXwRE&list=PLQh4-mYsu1HjzCaq-0ArETsOVf4Nqkzb0&index=2&t=0s)

