package parser

import (
	"strings"

	"github.com/rekey/go-club/common"
)

func parseTitle(title string) string {
	title = strings.ReplaceAll(title, " ", "-")
	return title
}

// Parse 根据URL解析对应的媒体信息，自动识别平台类型并调用相应的解析函数
// 参数u为媒体URL，返回值为解析后的媒体信息结构体指针，若平台不支持则返回nil
func Parse(u string) *common.Media {
	// pt parseType
	pt := common.GetUrlPT(u)
	switch pt {
	case common.PTKanav:
		result := parseKanav(u)
		result.Title = parseTitle(result.Title)
		return result
	case common.PTMadou:
		result := parseMadou(u)
		result.Title = parseTitle(result.Title)
		return result
	default:
		return nil
	}
}
