package spider

import "github.com/hu17889/go_spider/core/downloader"

//
//import (
//	"github.com/cnjack/novel-spider/model"
//	"log"
//	"time"
//)
//
//var MaxRunningTask = 50 //最大线程数
//var lock = make(chan int, MaxRunningTask)
//var StopSingle = false
//
var d = downloader.NewHttpDownloader()

type Novel struct {
	Title        string
	Auth         string
	Style        string
	Introduction string
	From         string
	Chapter      []*Chapter
}

type Chapter struct {
	Novel *Novel
	Title string
	From  string
	Data  string
}

//
//type Spider interface {
//	Name() string
//	Match(string) bool
//	Gain() (interface{}, error)
//}
//
//var spiders []Spider
//
//func Spider() {
//
//	for {
//		lock <- 0
//		go RunATask()
//		if StopSingle {
//			return
//		}
//	}
//}
//
//func RunATask() {
//	defer func() {
//		<-lock
//	}()
//	t, err := getTask()
//	if err != nil {
//		log.Println("get task err:", err)
//		return
//	}
//	if t == nil {
//		//如果没有任务就暂停一下
//		time.Sleep(10 * time.Second)
//	}
//	err = runTask(t)
//	if err != nil {
//		log.Println("runTask err", err)
//	}
//}
//
//func runTask(t *model.Task) error {
//	for _, s := range spiders {
//		if !s.Match(t.Url) {
//			continue
//		}
//		log.Println(s.Name(), " deal ", t.Url, " start")
//		resp, err := s.Gain()
//		if err != nil {
//			return err
//		}
//		//更新任务状态
//		err = flashTask(t, resp)
//		if err != nil {
//			return err
//		}
//		log.Println(s.Name(), " deal ", t.Url, " done")
//	}
//	return nil
//}
//
//func flashTask(t *model.Task, data interface{}) error {
//	return nil
//}
//
//func getTask() (task *model.Task, err error) {
//	db, err := model.MustGetDB()
//	if err != nil {
//		return nil, err
//	}
//	tx := db.Begin() //使用排他锁开启事务
//	if err = tx.Error; err != nil {
//		return nil, err
//	}
//	defer func() {
//		if err != nil {
//			tx.Rollback()
//		} else {
//			tx.Commit()
//		}
//	}()
//	task, err = model.FisrtTask(db)
//	if err != nil {
//		return nil, err
//	}
//	if task == nil {
//		return nil, nil
//	}
//	err = task.ChangeTaskStatus(model.TaskStatusRunning)
//	if err != nil {
//		return nil, err
//	}
//	return
//}
