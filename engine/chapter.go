package engine

import (
	"container/heap"
	"git.oschina.net/cnjack/go-task/task"
	"git.oschina.net/cnjack/novel-spider/config"
	"git.oschina.net/cnjack/novel-spider/mq"
	"log"
	"sync"
)

type ChaperTask struct {
	q *mq.Queue
	w *sync.WaitGroup
}

var ct *ChaperTask

func ChapterInit() *ChaperTask {
	ct := &ChaperTask{}
	ct.q = &mq.Queue{}
	return ct
}

func (c *ChaperTask) Push(task interface{}) {
	heap.Push(c.q, task)
}

func (c *ChaperTask) Pop() interface{} {
	return heap.Pop(c.q)
}

func (c *ChaperTask) Run() {
	c.w = &sync.WaitGroup{}
	for i := 0; i < config.GetSpiderConfig().MaxProcess; i++ {
		w.Add(1)
		c.Process()
	}
	w.Wait()
}

func (c *ChaperTask) Process() {
	defer func() {
		w.Done()
	}()
	taskI := c.Pop()
	task, ok := taskI.(*task.Task)
	if !ok {
		log.Printf("INFO: getTask error", taskI)
		c.Process()
	}
	err := runTask(task)
	if err != nil {
		//记录日志
		c.Push(task)
	}
	if config.GetSpiderConfig().StopSingle {
		return
	}
	c.Process()
}
