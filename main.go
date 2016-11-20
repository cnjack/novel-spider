package main

import (
	"fmt"

	"git.oschina.net/cnjack/novel-spider/engine"
	"git.oschina.net/cnjack/novel-spider/httpserv"
)

func main() {
	fmt.Println("serv running")
	go engine.Spider()

	httpserv.Http()
}
