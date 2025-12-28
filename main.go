package main

import (
	"log"
	"math"

	"github.com/rekey/go-club/downloader"
	"github.com/rekey/go-club/env"
	"github.com/rekey/go-club/parser"
)

const (
	ProgressDownloadComplete = 1.0 // 定义完成进度常量
	ProgressJoinIng          = 1.1
	ProgressComplete         = 100
)

func main() {
	// func() {
	// 	task := dao.NewTask("https://kanav.ad/index.php/vod/play/id/66923/sid/1/nid/1.html")
	// 	task.Save()
	// }()
	// func() {
	// 	task, err := dao.GetAll()
	// 	// task.SetStatusComplete()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	jsonb, err := json.Marshal(task)
	// 	log.Println(string(jsonb))
	// }()
	// 解析KANAV视频信息
	media := parser.Parse(
		"https://kanav.ad/index.php/vod/play/id/82550/sid/1/nid/1.html",
		// "https://madou.club/szl035-%e6%b7%ab%e6%ac%b2%e6%97%ba%e7%9b%9b%e7%9a%84%e5%a7%90%e5%a7%90%e8%89%b2%e8%af%b1%e8%a1%a8%e5%bc%9f%e5%ae%b6%e4%b8%ad%e5%81%9a%e7%88%b1.html",
	)
	if media == nil {
		log.Fatal("Failed to parse media information")
	}

	log.Printf("解析成功: 标题=%s, 制作人=%s, M3U8=%s", media.Title, media.Maker, media.M3u8Url)

	// // 开始下载
	downChan := downloader.Download(downloader.Args{
		Url:     media.M3u8Url,
		Referer: media.Url,
		Dir:     env.DownloadDir,
		Name:    media.Title + ".mp4",
		Num:     env.Concurrency,
	})

	// 监控下载进度
	for downResult := range downChan {
		switch downResult.Stat {
		case downloader.StatIng:
			p := downResult.Progress * 100
			if math.IsNaN(p) {
				p = 0
			}
			log.Printf("下载进度: %.1f%%", p)
		case downloader.StatJoin:
			log.Println("开始合并碎片")
		case downloader.StatComplete:
			log.Println("任务完成")
		}
	}

	// log.Println("程序结束")

	// log.Println("程序结束")
}
