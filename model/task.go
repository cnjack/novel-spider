package model

import "github.com/jinzhu/gorm"

type Task struct {
	gorm.Model
	Type        TaskType   `sql:"type"`
	Url         string     `sql:"url"`
	Status      TaskStatus `sql:"status"`
	Times       int        `sql:"times"`
	TargetID    int64      `sql:"target_id"`
	TargetField string     `sql:"targetfield"`
}

type TaskType int64

const (
	NovelTask TaskType = iota
	ChapterTask
)

type TaskStatus int64

const (
	TaskStatusPrepare TaskType = iota
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

func (t *Task) ChangeTaskStatus(tt TaskType) error {
	return db.Model(&Task{}).Where("id = ?", t.ID).Updates(map[string]interface{}{"status": tt}).Error
}
