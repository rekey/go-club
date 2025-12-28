package dao

import (
	"context"
	"log"

	"github.com/rekey/go-club/model"
	"gorm.io/gorm"
)

const (
	wait     = "wait"
	ing      = "ing"
	err      = "error"
	complete = "complete"
)

type Task struct {
	*model.Task
}

func initTask() {
	DB.AutoMigrate(&Task{})
	log.Println("init task")
}

func NewTask(url string) *Task {
	return &Task{
		Task: &model.Task{
			Url:    url,
			Status: wait,
		},
	}
}

func (t *Task) SetStatus(status string) {
	t.Status = status
	t.Save()
}

func (t *Task) SetStatusWait() {
	t.SetStatus(wait)
}

func (t *Task) SetStatusIng() {
	t.SetStatus(ing)
}

func (t *Task) SetStatusErr() {
	t.SetStatus(err)
}

func (t *Task) SetStatusComplete() {
	t.SetStatus(complete)
}
func (t *Task) Save() {
	DB.Save(t)
}

func GetAll() ([]Task, error) {
	return gorm.G[Task](DB).Find(context.Background())
}
