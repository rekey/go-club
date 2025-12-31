package dao

import (
	"log"

	"github.com/rekey/go-club/common"
	"github.com/rekey/go-club/model"
)

var TaskStats = struct {
	Wait     int
	Parser   int
	Download int
	Merge    int
	Move     int
	Err      int
	Complete int
}{
	Err:      -1,
	Wait:     0,
	Parser:   1,
	Download: 2,
	Merge:    3,
	Move:     4,
	Complete: 5,
}

type Task struct {
	*model.Task
}

func initTask() {
	DB.AutoMigrate(&Task{})
	log.Println("init task")
	resetStatus()
}

func resetStatus() {
	var tasks []Task
	err := DB.Model(Task{}).Where("status > ? and status < ?", TaskStats.Wait, TaskStats.Complete).Find(&tasks).Error
	if err != nil {
		log.Println("reset task status error", err)
		return
	}
	for _, task := range tasks {
		log.Println(task.Task)
		task.Status = TaskStats.Wait
		task.Progress = 0
		task.Save()
	}
}

// func resetStatus() {
// 	tasks, err := GetAllTask()
// 	if err != nil {
// 		log.Println("reset task status error", err)
// 		return
// 	}
// 	for _, task := range tasks {
// 		log.Println(task.Task)
// 		task.Status = TaskStats.Wait
// 		task.Progress = 0
// 		task.Save()
// 	}
// }

func NewTask(u string) *Task {
	return &Task{
		Task: &model.Task{
			Url:    u,
			Status: 0,
		},
	}
}

func (t *Task) UpdateStatus(status int) error {
	t.Status = status
	return t.Save()
}

func (t *Task) UpdateProgress(progress float64) error {
	t.Progress = progress
	return t.Save()
}

func (t *Task) Start() error {
	return t.UpdateStatus(TaskStats.Wait)
}

func (t *Task) StartDownload() error {
	return t.UpdateStatus(TaskStats.Download)
}

func (t *Task) SetName(name string) error {
	t.Name = name
	return t.Save()
}

func (t *Task) Save() error {
	return DB.Save(t).Error
}

func CreateTask(u string) *Task {
	pt := common.GetUrlPT(u)
	if pt == common.PTNone {
		return nil
	}
	task, err := FindTaskByURL(u)
	if task != nil && err == nil {
		return task
	}
	task = NewTask(u)
	task.Save()
	return task
}

func FindTaskByURL(url string) (*Task, error) {
	var task Task
	err := DB.Model(Task{}).Where("url = ?", url).First(&task).Error
	return &task, err
}

func UpdateTaskStatus(ur string, status int) error {
	task, err := FindTaskByURL(ur)
	if err != nil {
		return err
	}
	return task.UpdateStatus(1)
}

func GetAllTask() ([]Task, error) {
	var tasks []Task
	err := DB.Model(Task{}).Find(&tasks).Error
	return tasks, err
}

func GetOneWaitTask() (*Task, error) {
	var task Task
	err := DB.Model(Task{}).Where("status = ?", 0).First(&task).Error
	return &task, err
}
