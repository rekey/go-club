package common

import (
	"log"
	"net/http"
	"net/url"
	"os"
)

type Media struct {
	Title   string
	Maker   string
	Tags    []string
	Url     string
	M3u8Url string
}

type ParseType int

const (
	PTMadou ParseType = iota
	PTKanav
)

// IsValid 检查 ParseType 是否在有效范围内
func (parseType ParseType) IsValid() bool {
	return parseType >= PTMadou && parseType <= PTKanav
}

var hostPTMap = map[string]ParseType{
	"madou.club": PTMadou,
	"kanav.ad":   PTKanav,
}

var referMap = map[ParseType]string{
	PTMadou: "https://madou.club",
	PTKanav: "https://kanav.ad",
}

func GetUrlPT(u string) ParseType {
	uo, err := url.Parse(u)
	if err != nil {
		return -1
	}
	if value, ok := hostPTMap[uo.Host]; ok {
		return value
	}
	return -1
}

func GetUrlRefer(u string) string {
	pt := GetUrlPT(u)
	if value, ok := referMap[pt]; ok {
		return value
	}
	return ""
}

// GetUrlType 返回给定URL的类型，通过调用GetUrlRefer实现
func GetUrlType(u string) string {
	return GetUrlRefer(u)
}

func GetHeader(url string) map[string]string {
	return map[string]string{
		"accept":             "*/*",
		"accept-language":    "zh-CN,zh;q=0.5",
		"priority":           "u=1, i",
		"sec-ch-ua":          "\"Brave\";v=\"143\", \"Chromium\";v=\"143\", \"Not A(Brand\";v=\"24\"",
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": "\"Linux\"",
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "cross-site",
		"sec-gpc":            "1",
		"Referer":            url,
	}
}

func GetHttpHeader(url string) http.Header {
	header := GetHeader(url)
	result := http.Header{}
	for k, v := range header {
		result.Set(k, v)
	}
	return result
}

// httpGet 发送HTTP GET请求并返回响应
func HttpGet(u string, refer string) (*http.Response, error) {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	// 设置请求头
	headerRefer := refer
	if headerRefer == "" {
		headerRefer = u
	}
	req.Header = GetHttpHeader(headerRefer)
	return http.DefaultClient.Do(req)
}

func dirExists(path string) bool {
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		return true
	}
	return false
}

func CreateDir(path string) {
	if dirExists(path) {
		log.Println("Directory already exists:", path)
		return
	}
	err := os.RemoveAll(path)
	if err != nil {
		log.Println("Failed to remove directory:", path, err)
		return
	}
	err = os.Mkdir(path, 0755)
	if err != nil {
		log.Println("Failed to create directory:", path, err)
		return
	}
	log.Println("Directory created:", path)
}
