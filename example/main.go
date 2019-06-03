package main

import (
	"fmt"
	"github.com/kuangwanjing/ruleparser/parser"
	"strconv"
	"strings"
)

type SoftwareInfo struct {
	Sid     string  `rule:"sid"`
	Ver     Version `rule:"ver"`
	Channel string  `rule:"channel"`
	Count   int     `rule:"cnt"`
	F       bool    `rule:"f"`
}

type Version struct {
	value string
}

func (ver Version) Cmp(val string) (int, error) {
	vs1 := strings.Split(ver.value, ".")
	vs2 := strings.Split(val, ".")
	c1 := 0
	c2 := 0

	for c1 < len(vs1) && c2 < len(vs2) {
		v1, err1 := strconv.Atoi(vs1[c1])
		v2, err2 := strconv.Atoi(vs2[c2])

		if err1 != nil {
			return -1, err1
		}

		if err2 != nil {
			return -1, err2
		}

		if v1 < v2 {
			return -1, nil
		} else if v1 > v2 {
			return 1, nil
		}

		c1 += 1
		c2 += 1
	}

	if c1 != len(vs1) {
		return 1, nil
	} else if c2 != len(vs2) {
		return -1, nil
	}

	return 0, nil
}

func (ver Version) In(val string) (int, error) {
	vs := strings.Split(val, ",")

	for _, v := range vs {
		if v == ver.value {
			return 0, nil
		}
	}

	return -1, nil
}

func main() {
	rules := "ver < `3.5.0`;ver > `1.5.0`;ver in `2.5.0,2.5.1`;channel==`google play`;cnt >= -2;f==0"
	software := SoftwareInfo{"134efa", &Version{"2.5.0"}, "google play", 3, false}
	p, err := parser.ParserInit(rules)

	fmt.Println(p)
	fmt.Println(err)

	if err == nil {
		fmt.Println(p.Examine(&software))
		fmt.Println(p.Examine(software))
		fmt.Println(p.Examine([]string{"abc"}))
		fmt.Println(p.Examine(float32(13.22)))
	}
}
