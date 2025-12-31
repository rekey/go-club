package parser

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/url"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/rekey/go-club/common"
)

type tKanavData struct {
	Data struct {
		Name  string `json:"vod_name"`
		Maker string `json:"vod_director"`
		Actor string `json:"vod_actor"`
	} `json:"vod_data"`
	M3U8 string `json:"url"`
}

func parseKanavData(s string) *tKanavData {
	re := regexp.MustCompile(`var\s+player_aaaa\s*=\s*(\{.*\})`)
	matches := re.FindStringSubmatch(s)
	if len(matches) < 2 {
		log.Println("No JSON data found in script")
		return nil
	}

	jsonStr := matches[1]
	// fmt.Printf("提取的JSON: %s\n", jsonStr) // 调试输出

	// 解析JSON
	var result tKanavData
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		log.Printf("Failed to parse JSON: %v", err)
		return nil
	}
	m3u8UrlBit, _ := base64.URLEncoding.DecodeString(result.M3U8)
	m3u8Url, _ := url.QueryUnescape(string(m3u8UrlBit))
	result.M3U8 = string(m3u8Url)
	return &result
}

func parseKanav(u string) *Media {
	res, err := common.HttpGet(u, common.GetUrlRefer(u))
	if err != nil {
		log.Println(err)
		return nil
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println(err)
		return nil
	}
	s := doc.Find(".video-server .poplayer script").First().Text()
	result := parseKanavData(s)
	var media = Media{
		Url:     u,
		Title:   result.Data.Name,
		Maker:   result.Data.Maker,
		M3u8Url: result.M3U8,
	}
	media.Tags = doc.Find(".video-countext-tags a").Map(func(i int, s *goquery.Selection) string {
		return s.Text()
	})
	return &media
}
