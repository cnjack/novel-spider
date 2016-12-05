package httpserv

import (
	"time"

	"git.oschina.net/cnjack/novel-spider/model"
)

type Status struct {
	ServerStartTime   string `json:"server_start_time"`
	ServerRunningTime string `json:"server_running_time"`
	NovelNum          int    `json:"novel_num"`
}

var status = Status{}

func ReloadStatus() {
	go func() {
		for range time.Tick(1 * time.Second) {
			status.ServerRunningTime = time.Now().Sub(startTime).String()
		}
	}()
	go func() {
		for range time.Tick(60 * time.Second) {
			status.NovelNum, _ = model.CountNovel()
		}
	}()
	status.ServerStartTime = startTime.Format(time.RFC3339)
	status.NovelNum, _ = model.CountNovel()
}
