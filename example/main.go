package main

import (
	"fmt"
	"parser"
)

type SoftwareInfo struct {
	Sid     string  `rule:"sid"`
	Ver     Version `rule:"ver"`
	Channel string  `rule:"channel"`
}

type Version struct {
	value string
}

func main() {
	rules := "ver<`3.0.0`;ver>`1.5.0`;channel==`google play`"
	software := SoftwareInfo{"134efa", Version{"3.5.0"}, "google play"}
	p, err := parser.ParserInit(rules)

	fmt.Println(p)

	if err == nil {
		fmt.Println(p.Examine(software))
	}
}
