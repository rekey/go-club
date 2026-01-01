package parser

import (
	"strings"

	"github.com/rekey/go-club/common"
)

type Media struct {
	Title   string
	Maker   string
	Actor   string
	Tags    []string
	Url     string
	M3u8Url string
	Thumb   string
}

func sanitize(s string) string {
	// 1. 使用标准库的TrimSpace处理所有Unicode空白字符
	//    这会移除字符串开头和结尾的所有空格、Tab、换行等空白字符
	s = strings.TrimSpace(s)

	if s == "" {
		return "unknown"
	}

	// 2. 先替换路径分隔符（这是核心安全措施）
	//    这样后续的Trim操作才不会因为中间有分隔符而误处理
	s = strings.ReplaceAll(s, "/", "-")
	s = strings.ReplaceAll(s, "\\", "-")

	// 3. 移除字符串首尾的特殊字符（修正原代码的顺序问题）
	//    原代码多个Trim调用，实际上只移除第一个匹配字符，改为循环Trim
	trimChars := "/\\-"
	for {
		trimmed := strings.Trim(s, trimChars)
		if trimmed == s {
			break // 没有更多字符可移除
		}
		s = trimmed
	}

	// 4. 处理字符串内部的空白字符（包括Tab、换行等）
	//    使用strings.Fields分割再合并，可处理所有Unicode空白字符
	//    先分割再合并，比正则表达式更高效
	parts := strings.Fields(s) // Fields按空白分割，包含空格、Tab、换行等
	if len(parts) == 0 {
		return "unknown"
	}
	s = strings.Join(parts, "_") // 用下划线连接非空白部分

	// 5. 移除可能产生的连续分隔符
	//    例如 "a---b" -> "a-b"，"--a-b--" -> "a-b"
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}

	// 6. 最终检查：如果经过处理又变成空字符串，返回默认值
	if s == "" {
		return "unknown"
	}

	return s
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
		result.Title = sanitize(result.Title)
		result.Maker = sanitize(result.Maker)
	}
	if result.Maker == "" {
		result.Maker = "unknown"
	}
	return result
}
