package main

import (
	"fmt"
	"parser"
	"strconv"
	"strings"
)

type SoftwareInfo struct {
	Sid     string  `rule:"sid"`
	Ver     Version `rule:"ver"`
	Channel string  `rule:"channel"`
}

type Version struct {
	value string
}

func (ver Version) Cmp(val string) int {
	vs1 := strings.Split(ver.value, ".")
	vs2 := strings.Split(val, ".")
	c1 := 0
	c2 := 0

	for c1 < len(vs1) && c2 < len(vs2) {
		v1, _ := strconv.Atoi(vs1[c1])
		v2, _ := strconv.Atoi(vs2[c2])
		if v1 < v2 {
			return -1
		} else if v1 > v2 {
			return 1
		}

		c1 += 1
		c2 += 1
	}

	if c1 != len(vs1) {
		return 1
	} else if c2 != len(vs2) {
		return -1
	}

	return 0
}

func main() {
	rules := "ver<`3.0.0`;ver>`1.5.0`;channel==`google play`"
	software := SoftwareInfo{"134efa", Version{"2.5.0"}, "google play"}
	p, err := parser.ParserInit(rules)

	fmt.Println(p)

	if err == nil {
		fmt.Println(p.Examine(software))
	}
}
