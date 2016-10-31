package main

import (
	"encoding/json"
	"fmt"
	"github.com/cnjack/novel-spider/spider"
)

func main() {
	s := &spider.Snwx{}
	s.Match("http://www.snwx.com/book/124/124785/")
	resp, err := s.Gain()
	if err != nil {
		fmt.Println(err)
	}
	jsonByte, err := json.Marshal(resp)
	fmt.Println(string(jsonByte))
}
