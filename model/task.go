package model

import "github.com/jinzhu/gorm"

type Task struct {
	gorm.Model
	TType       TaskType   `sql:"ttype"`
	Url         string     `sql:"url"`
	Status      TaskStatus `sql:"status"`
	Times       int        `sql:"times"`
	TargetID    uint       `sql:"target_id"`
	TargetField string     `sql:"targetfield"`
}

type TaskType uint8

const (
	NovelTask TaskType = iota
	ChapterTask
)

type TaskStatus uint8

const (
	TaskStatusPrepare TaskStatus = iota
	TaskStatusRunning
	TaskStatusFail
	TaskStatusOk
)

func FisrtTask(db *gorm.DB) (*Task, error) {
	t := &Task{}
	if err := db.Where("status in (?, ?)", TaskStatusPrepare, TaskStatusFail).Order("id desc").First(t).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return t, nil
}

func CountTasks(t ...TaskStatus) (count int, err error) {
	err = db.Model(&Task{}).Where("status in (?)", t).Count(&count).Error
	return
}

func (t *Task) ChangeTaskStatus(tt TaskStatus) error {
	return db.Model(&Task{}).Where("id = ?", t.ID).Updates(map[string]interface{}{"status": tt}).Error
}
