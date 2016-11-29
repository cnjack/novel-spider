package main

import (
	"fmt"

	"git.oschina.net/cnjack/novel-spider/httpserv"
	"git.oschina.net/cnjack/novel-spider/job"
)

func main() {
	fmt.Println("serv running")
	go job.Spider()

	httpserv.Http()
}
