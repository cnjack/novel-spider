package main

import (
	"fmt"

	"git.oschina.net/cnjack/novel-spider/spider"
)

func main() {
	s := &spider.SnwxChapter{}
	s.Match("http://www.snwx.com/book/124/124785/26045253.html")
	resp, err := s.Gain()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp)
}
