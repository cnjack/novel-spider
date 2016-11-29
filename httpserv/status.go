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

func appendStatus(ss []Status, s Status) []Status {
	if len(ss) >= 40 {
		sss := make([]Status, len(ss))
		for k, v := range ss {
			if k == 0 {
				continue
			}
			sss[k-1] = v
		}
		sss[len(ss)-1] = s
		return sss
	}
	return append(ss, s)
}

var statuss = []Status{}

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
	go func() {
		for range time.Tick(5 * time.Second) {
			statuss = appendStatus(statuss, status)
		}
	}()
	status.ServerStartTime = startTime.Format(time.RFC3339)
	status.ServerTunningTime = time.Now().Sub(startTime).String()
	status.RunningTasks, _ = model.CountTasks(model.TaskStatusRunning)
	status.PrepareTasks, _ = model.CountTasks(model.TaskStatusPrepare, model.TaskStatusFail)
	status.SuccessTasks, _ = model.CountTasks(model.TaskStatusOk)
	status.NovelNum, _ = model.CountNovel()
	statuss = appendStatus(statuss, status)
}
