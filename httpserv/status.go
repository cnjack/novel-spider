package httpserv

import (
	"time"

	"git.oschina.net/cnjack/novel-spider/model"
)

type Status struct {
	RunningTasks      int    `json:"running_tasks"`
	PrepareTasks      int    `json:"prepare_tasks"`
	ServerStartTime   string `json:"server_start_time"`
	ServerTunningTime string `json:"server_running_time"`
	Now               string `json:"now"`
	SuccessTasks      int    `json:"success_tasks"`
	NovelNum          int    `json:"novel_num"`
}

var status = Status{}

func ReloadStatus() {
	go func() {
		for range time.Tick(1 * time.Second) {
			status.ServerTunningTime = time.Now().Sub(startTime).String()
			status.Now = time.Now().Format("15:04:05")
		}
	}()
	go func() {
		for range time.Tick(15 * time.Second) {
			status.RunningTasks, _ = model.CountTasks(model.TaskStatusRunning)
			status.PrepareTasks, _ = model.CountTasks(model.TaskStatusPrepare, model.TaskStatusFail)
			status.SuccessTasks, _ = model.CountTasks(model.TaskStatusOk)
		}
	}()
	go func() {
		for range time.Tick(60 * time.Second) {
			status.NovelNum, _ = model.CountNovel()
		}
	}()
	status.ServerStartTime = startTime.Format(time.RFC3339)
	status.ServerTunningTime = time.Now().Sub(startTime).String()
	status.RunningTasks, _ = model.CountTasks(model.TaskStatusRunning)
	status.PrepareTasks, _ = model.CountTasks(model.TaskStatusPrepare, model.TaskStatusFail)
	status.SuccessTasks, _ = model.CountTasks(model.TaskStatusOk)
	status.NovelNum, _ = model.CountNovel()
}
