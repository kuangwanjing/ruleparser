// This example shows a realistic scenario when applications get data in json format and filter the data by its need.
package main

import (
	"encoding/json"
	"fmt"
	"github.com/kuangwanjing/ruleparser/parser"
)

// The struct video is exposed to rule parser with 3 fields: uploader, category
type Video struct {
	ID       int       `json:"id" rule:"-"`
	Uploader string    `json:"uploader" rule:"uploader"`
	VTags    VideoTags `json:"tags" rule:"tags"` // add customize operation on tags
}

// use another struct as a receiver of the operation method instead of an array of string
type VideoTags struct {
	Tags []string
}

/*
provide customized marshaling and unmarshaling behaviors for the json encoding and decoding procedures to fix the difference between struct tags and data structure
*/
func (vtags *VideoTags) UnmarshalJSON(b []byte) error {
	var tags []string
	if err := json.Unmarshal(b, &tags); err != nil {
		return err
	}
	vtags.Tags = tags
	return nil
}

func (vtags *VideoTags) MarshalJSON() ([]byte, error) {
	return json.Marshal(vtags.Tags)
}

// define operation `has` on VideoTags
func (vtags VideoTags) Has(pattern string) (int, error) {
	for _, tag := range vtags.Tags {
		if tag == pattern {
			return 0, nil // if tag is found, return 0 as a successful operation
		}
	}
	return -1, nil // return non-zero as a failed operation
}

func main() {

	// in this example, we have three video objects to be filtered by the rules
	blob := `[{"id":1,"uploader":"user_1","tags":["sports"]},{"id":2,"uploader":"user_2","tags":["pets","live"]},{"id":3,"uploader":"user_3","tags":["sports","live"]}]`

	// unmarshal the blob into objects
	var videos []Video
	if err := json.Unmarshal([]byte(blob), &videos); err != nil {
		fmt.Println(err)
	}

	fmt.Println(videos)

	// set the rules and get a rule parser
	rules := "uploader != `uploader_2`;tags has `sports`"
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

		// print the videos matches with the rules
		fmt.Println(filteredData)
	}
}
