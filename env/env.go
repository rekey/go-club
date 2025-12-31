package env

import (
	"log"
	"os"
	"path"
	"path/filepath"

	"go.yaml.in/yaml/v4"
)

type Config struct {
	DataDir     string `yaml:"data_dir"`
	DownloadDir string `yaml:"download_dir"`
	LogDir      string `yaml:"log_dir"`
	Concurrency int    `yaml:"concurrency"`
	TaskNum     int    `yaml:"task_num"`
}

var configFile = "./config.yaml"
var DataDir = "./data"
var DownloadDir = "./download"
var Concurrency = 20
var TaskNum = 5
var LogDir = "./logs"

// init 初始化应用配置，从配置文件中读取并解析 YAML 格式的配置项
// 自动处理相对路径转换为绝对路径，并设置以下全局变量：
// - DataDir: 数据存储目录
// - DownloadDir: 下载文件存储目录
// - Concurrency: 并发下载数
// - TaskNum: 任务数量限制
// 如果配置文件读取或解析失败，会触发 panic 终止程序
func init() {
	cwd, _ := os.Getwd()
	configFile = path.Join(cwd, configFile)
	configContent, err := os.ReadFile(configFile)
	if err != nil {
		log.Println("read config file error:", err)
		panic(err)
	}
	var config Config
	err = yaml.Unmarshal(configContent, &config)
	if err != nil {
		log.Println("unmarshal config file error:", err)
		panic(err)
	}
	DataDir = config.DataDir
	if !path.IsAbs(config.DataDir) {
		DataDir = path.Join(cwd, config.DataDir)
	}
	DataDir, _ = filepath.Abs(DataDir)
	DownloadDir = config.DownloadDir
	if !path.IsAbs(config.DownloadDir) {
		DownloadDir = path.Join(cwd, config.DownloadDir)
	}
	DownloadDir, _ = filepath.Abs(DownloadDir)
	LogDir = config.LogDir
	if !path.IsAbs(config.LogDir) {
		LogDir = path.Join(cwd, config.LogDir)
	}
	LogDir, _ = filepath.Abs(LogDir)
	if config.Concurrency > 0 {
		Concurrency = config.Concurrency
	}
	if config.TaskNum > 0 {
		TaskNum = config.TaskNum
	}
	log.Println("参数名", "配置文件", "最终使用")
	log.Println("DataDir", config.DataDir, DataDir)
	log.Println("DownloadDir", config.DownloadDir, DownloadDir)
	log.Println("LogDir", config.LogDir, LogDir)
	log.Println("Concurrency", config.Concurrency, Concurrency)
	log.Println("TaskNum", config.TaskNum, TaskNum)
}
