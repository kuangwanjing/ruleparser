// This example shows how to use ruleparser to filter data with basic data fields
package main

import (
	"fmt"
	"github.com/kuangwanjing/ruleparser/parser"
)

// The struct video is exposed to rule parser with 3 fields: uploader, category
type Video struct {
	ID         int    `rule:"-"` // adding any rules on ID won't make any difference.
	Uploader   string `rule:"uploader"`
	Category   string `rule:"category"`
	UploadTime int64  `rule:"ut"`
}

func main() {

	// in this example, we have three video objects to be filtered by the rules
	videos := []Video{
		Video{1, "uploader_1", "sports", 1546300800}, // not match with upload_time
		Video{2, "uploader_2", "pets", 1546646400},   // not match with the category
		Video{3, "uploader_3", "sports", 1546473600},
	}

	// set the rules and get a rule parser
	rules := "uploader != `uploader_2`;category == `sports`;ut >= 1546387200"
	p, err := parser.ParserInit(rules)

	var filteredData []Video

	if err != nil {
		fmt.Println(err)
	} else {
		for _, v := range videos {
			// examine each video with the above three rules by the parser
			rst, err := p.Examine(&v)

			// if there's no error happens and the video passes the examination
			if err == nil && rst {
				filteredData = append(filteredData, v)
			}
		}

		fmt.Println(filteredData)
	}
}
