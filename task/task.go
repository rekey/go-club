package task

import (
	"errors"
	"log"
	"math"
	"sync"
	"time"

	"github.com/rekey/go-club/dao"
	"github.com/rekey/go-club/downloader"
	"github.com/rekey/go-club/env"
	"github.com/rekey/go-club/parser"
	"github.com/rekey/go-club/pool"
)

var Stats = dao.TaskStats

type Task struct {
	Task       *dao.Task
	Media      *parser.Media
	Downloader downloader.Result
}

func parse(task *Task) error {
	task.Task.UpdateStatus(Stats.Parser)
	task.Media = parser.Parse(task.Task.Url)
	if task.Media == nil {
		task.Task.Status = Stats.Err
		task.Task.Err = "parse error"
		return errors.New("parse error")
	}
	task.Task.SetName(task.Media.Title)
	return nil
}

func download(task *Task) error {
	task.Task.UpdateStatus(Stats.Download)
	name := task.Media.Title + ".mp4"
	downChan := downloader.Download(downloader.Args{
		Url:     task.Media.M3u8Url,
		Referer: task.Task.Url,
		Dir:     env.DownloadDir,
		Name:    name,
		Num:     env.Concurrency,
	})
	var err error = nil
	// 监控下载进度
	for downResult := range downChan {
		switch downResult.Stat {
		case downloader.StatIng:
			p := downResult.Progress * 100
			if math.IsNaN(p) {
				p = 0
			}
			task.Task.UpdateProgress(p)
			log.Println(name, "下载进度", p, "%")
		case downloader.StatJoin:
			task.Task.UpdateProgress(100)
			task.Task.UpdateStatus(Stats.Merge)
			log.Println(name, "开始合并碎片")
		case downloader.StatComplete:
			task.Task.UpdateProgress(100)
			task.Task.UpdateStatus(Stats.Complete)
			log.Println(name, "任务完成")
		case downloader.StatErr:
			task.Task.UpdateProgress(0)
			task.Task.UpdateStatus(Stats.Err)
			err = errors.New("download error")
			log.Println(name, "任务出错")
		}
	}
	if err == nil {
		task.Task.UpdateStatus(Stats.Complete)
	}
	return err
}

var getTaskMutex sync.RWMutex

func runTask(id int) {
	getTaskMutex.Lock()
	daoTask, err := dao.GetOneWaitTask()
	if err != nil || daoTask == nil {
		getTaskMutex.Unlock()
		// log.Println("get one wait task error", err)
		time.Sleep(time.Second)
		return
	}
	daoTask.StartDownload()
	getTaskMutex.Unlock()
	log.Println("task start", "id", id, "url", daoTask.Url)
	task := &Task{
		Task: daoTask,
	}
	err = parse(task)
	if err != nil {
		log.Println("parse error", task.Task.Url, err)
		return
	}
	download(task)
}

func Run() {
	log.Println("任务模块", "启动", "并发任务数量", env.TaskNum, "任务下载并发数量", env.Concurrency)
	p := pool.NewPool(env.TaskNum, runTask, true)
	p.Run()
}
