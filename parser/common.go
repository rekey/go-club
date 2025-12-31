package parser

import (
	"strings"

	"github.com/rekey/go-club/common"
)

type Media struct {
	Title   string
	Maker   string
	Tags    []string
	Url     string
	M3u8Url string
}

func parseTitle(title string) string {
	title = strings.ReplaceAll(title, " ", "_")
	return title
}

// Parse 根据URL解析对应的媒体信息，自动识别平台类型并调用相应的解析函数
// 参数u为媒体URL，返回值为解析后的媒体信息结构体指针，若平台不支持则返回nil
func Parse(u string) *Media {
	// pt parseType
	pt := common.GetUrlPT(u)
	var result *Media
	switch pt {
	case common.PTKanav:
		result = parseKanav(u)
	case common.PTMadou:
		result = parseMadou(u)
	}
	if result != nil {
		result.Title = parseTitle(result.Title)
	}
	return result
}
