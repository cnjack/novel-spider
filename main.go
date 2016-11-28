package main

import (
	"fmt"

	"git.oschina.net/cnjack/novel-spider/job"
	"git.oschina.net/cnjack/novel-spider/httpserv"
)

func main() {
	fmt.Println("serv running")
	go job.Spider()

	httpserv.Http()
}
