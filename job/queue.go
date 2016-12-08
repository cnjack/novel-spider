package job

import (
	"encoding/json"
	"log"
	"time"

	"git.oschina.net/cnjack/novel-spider/config"
	"git.oschina.net/cnjack/novel-spider/model"
	"git.oschina.net/cnjack/novel-spider/mq"
	"github.com/jinzhu/gorm"
)

var b *mq.Broker

func init() {
	b = &mq.Broker{
		Addrs: config.GetSpiderConfig().BrokerAddr,
	}
	for i:=0; i<config.GetSpiderConfig().MaxProcess;i++ {
		if _, err := b.Subscribe(config.GetSpiderConfig().BrokerTopic, config.GetSpiderConfig().BrokerChannel, nTaskHandle); err != nil {
			panic(err)
		}
	}
	if err := b.Connect(); err != nil {
		panic(err)
	}
}

func PublishTask(task *model.Task) (err error) {
	taskString, _ := json.Marshal(task)
	err = b.Publish(config.GetSpiderConfig().BrokerTopic, &mq.Message{
		Header: nil,
		Body:   []byte(taskString),
	})
	return
}

func nTaskHandle(p mq.Publication) (err error) {
	body := p.Message().Body
	var task model.Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		//记录日志
		log.Println("INFO: get task error", err)
	}
	err = RunTask(&task)
	if err != nil {
		//记录日志
		log.Println("INFO: runTask error", err)
	}
	err = p.Ack()
	if err != nil {
		log.Println("INFO: runTask mark success error", err)
	}
	return nil
}

func UpdateNovelTask() {
	task, err := getTask()
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Printf("ERROR: getTask ERROR; ERRstring: %s", err.Error())
		}
		time.Sleep(30 * time.Second)
		UpdateNovelTask()
		return
	}
	log.Printf("INFO: getTask ok; task id: %d", task.ID)
	taskString, _ := json.Marshal(task)
	err = b.Publish(config.GetSpiderConfig().BrokerTopic, &mq.Message{
		Header: nil,
		Body:   []byte(taskString),
	})
	if err != nil {
		log.Println("INFO: publish task err:", err)
	}
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
