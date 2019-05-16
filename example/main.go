package main

import (
	"fmt"
	"parser"
)

func main() {
	rules := "ver<`3.0.0`;ver>`1.5.0`;channel==`google play`"
	fmt.Println(parser.ParserInit(rules))
}
