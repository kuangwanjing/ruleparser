// define the comparison behavior between software versions
package main

import (
	"fmt"
	"github.com/kuangwanjing/ruleparser/parser"
	"strconv"
	"strings"
)

type SoftwareInfo struct {
	Sid     string   `rule:"-"`
	Ver     *Version `rule:"ver"`
	Channel string   `rule:"channel"`
	Count   int      `rule:"cnt"`
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

		// if the current version is smaller, return a negative integer
		if v1 < v2 {
			return -1, nil
		} else if v1 > v2 {
			// if the current version is greater, return a positive integer
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

	// otherwise two versions are the same, return 0
	return 0, nil
}

func main() {
	rules := "ver < `3.5.0`;ver > `1.5.0`;channel==`google play`;cnt >= -2"
	software := SoftwareInfo{"134efa", &Version{"2.5.0"}, "google play", 3}
	p, err := parser.ParserInit(rules)

	if err == nil {
		fmt.Println(p.Examine(&software))
		fmt.Println(p.Examine(software))
	}
}
