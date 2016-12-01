package job

import (
	"log"
	"sync"
	"time"

	"git.oschina.net/cnjack/novel-spider/config"
	"git.oschina.net/cnjack/novel-spider/model"
	"git.oschina.net/cnjack/novel-spider/mq"
	"github.com/jinzhu/gorm"
)

var q = mq.New(0)

type task struct {
}

var wg = sync.WaitGroup{}

func Run() {
	var t task
	for i := 0; i < config.GetSpiderConfig().MaxProcess; i++ {
		wg.Add(1)
		go t.Process()
	}
	wg.Wait()
}

func (t task) Process() {
	defer func(){
		wg.Done()
	}()
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

func UpdateNovelTask() {
	t, err := getTask()
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Printf("ERROR: getTask ERROR; ERRstring: %s", err.Error())
		}
		time.Sleep(30 * time.Second)
		UpdateNovelTask()
	}
	log.Printf("INFO: getTask ok; task id: %d", t.ID)
	q.PutNoWait(t)
	if config.GetSpiderConfig().StopSingle {
		return
	}
	UpdateNovelTask()
}

func getTask() (task *model.Task, err error) {
	db, err := model.MustGetDB()
	if err != nil {
		return nil, err
	}
	tx := db.Begin()
	if err = tx.Error; err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	task = &model.Task{}
	//如果准备中则24小时候启动 如果是错误则直接启动
	if err = db.Raw("SELECT * FROM tasks WHERE (updated_at > ? AND status = ?) OR status = ? AND deleted_at IS NULL ORDER BY status, id ASC LIMIT 0, 1 FOR UPDATE", time.Now().Add(24*time.Hour), model.TaskStatusPrepare, model.TaskStatusFail).Scan(task).Error; err != nil {
		return nil, err
	}
	if err = tx.Model(&model.Task{}).Where("id = ?", task.ID).Updates(map[string]interface{}{"status": model.TaskStatusRunning, "updated_at": time.Now()}).Error; err != nil {
		return nil, err
	}
	return
}
