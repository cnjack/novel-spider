package engine

import (
	"log"
	"sync"
	"time"

	"git.oschina.net/cnjack/novel-spider/config"
	"git.oschina.net/cnjack/novel-spider/model"
	"git.oschina.net/cnjack/novel-spider/mq"
)

var q = mq.New(0)

type task struct {
}

func Run() {
	var wg = sync.WaitGroup{}
	var t task
	for i := 0; i < config.GetSpiderConfig().MaxProcess; i++ {
		wg.Add(1)
		go t.Process()
	}
	wg.Wait()
}

func (t task) Process() {
	if q.IsEmpty() {
		time.Sleep(1 * time.Second)
		t.Process()
	}
	taskI, err := q.GetNoWait()
	if err != nil {
		log.Println("INFO: getTask error:", err)
		t.Process()
	}
	task, ok := taskI.(*model.Task)
	if !ok {
		log.Println("INFO: getTask error not right")
		t.Process()
	}
	err = RunTask(task)
	if err != nil {
		//记录日志
		log.Println("INFO: runTask error", err)
		q.PutNoWait(task)
	}
	if config.GetSpiderConfig().StopSingle {
		return
	}
	t.Process()
}
