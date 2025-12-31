package downloader

import (
	"os"
	"path"
	"time"

	"github.com/canhlinh/hlsdl"
	"github.com/rekey/go-club/common"
	"github.com/rekey/go-club/env"
)

const (
	StatIng      = "ing"
	StatErr      = "error"
	StatJoin     = "join"
	StatComplete = "complete"
)

/*
init 初始化下载目录，确保下载目录存在
*/
func init() {
	DownloadDir := env.DownloadDir
	common.CreateDir(DownloadDir)
}

type Args struct {
	Url     string
	Referer string
	Dir     string
	Name    string
	Num     int
}

type Result struct {
	Progress float64
	Stat     string
}

var inited = false

func clear(dir string) {
	if inited {
		return
	}
	inited = true
	tmp := path.Join(dir, "tmp")
	common.CreateDir(tmp)
}

/*
Download 启动异步HLS下载任务并返回进度通知通道

参数:

	args: 下载参数，包含URL、引用头、目录等配置

返回:

	chan Result: 进度通知通道，会发送以下状态:
	  - StatIng: 下载中，包含当前进度(0-1)
	  - StatJoin: 下载完成合并中
	  - StatComplete: 下载完成
	  - StatErr: 下载出错

注意:

	通道会在下载结束后自动关闭，调用方应处理所有状态通知
*/
func Download(args Args) chan Result {
	clear(args.Dir)
	tmp := path.Join(args.Dir, "tmp")
	tmpFile := path.Join(tmp, args.Name)
	originFile := path.Join(args.Dir, args.Name)
	hlsDL := hlsdl.New(args.Url, common.GetHeader(args.Referer), tmp, args.Name, args.Num, false)
	progressChan := make(chan Result)
	progressChanClose := false
	go func() {
		defer func() {
			progressChanClose = true
			close(progressChan)
		}()
		_, err := hlsDL.Download()
		if err != nil {
			progressChan <- Result{
				Progress: 0,
				Stat:     StatErr,
			}
			return
		}
		os.Rename(tmpFile, originFile)
		progressChan <- Result{
			Progress: 1,
			Stat:     StatComplete,
		}
	}()
	go func() {
		for {
			if progressChanClose {
				break
			}
			p := hlsDL.GetProgress()
			if p >= 1 {
				progressChan <- Result{
					Progress: p,
					Stat:     StatJoin,
				}
				break
			}
			progressChan <- Result{
				Progress: p,
				Stat:     StatIng,
			}
			time.Sleep(time.Millisecond * 1000)
		}
	}()
	return progressChan
}
