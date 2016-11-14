package main

import (
	"fmt"

	"git.oschina.net/cnjack/novel-spider/engine"
)

func main() {
	fmt.Println("serv running")
	go engine.Spider()

	engine.Http()
}
