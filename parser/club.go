package parser

import (
	"io"
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/rekey/go-club/common"
)

func parseMadouM3U8Url(u string) string {
	if u == "" {
		return ""
	}
	res, err := common.HttpGet(u, common.GetUrlRefer(u))
	if err != nil {
		log.Println(err)
		return ""
	}
	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return ""
	}
	re := regexp.MustCompile(`var\s+(token|m3u8)\s*=\s*['"]([^'"]+)['"]`)
	matches := re.FindAllStringSubmatch(string(bodyBytes), -1)
	result := make(map[string]string)
	for _, match := range matches {
		result[match[1]] = match[2]
	}
	token := result["token"]
	m3u8 := result["m3u8"]
	m3u8UrlOrigin := "https://dash.madou.club" + m3u8 + "?token=" + token
	return m3u8UrlOrigin
}

func parseMadou(u string) *Media {
	res, err := common.HttpGet(u, common.GetUrlRefer(u))
	if err != nil {
		log.Println(err)
		return nil
	}
	defer res.Body.Close()
	media := Media{
		Url: u,
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println(err)
		return nil
	}
	media.Title = doc.Find(".article-title").First().Text()
	media.Maker = doc.Find(".article-meta .item-3").First().Text()
	media.Tags = doc.Find(".article-tags a").Map(func(i int, s *goquery.Selection) string {
		return s.Text()
	})
	media.M3u8Url, _ = doc.Find(".article-content iframe").First().Attr("src")
	// data, _ := json.Marshal(media)
	// log.Println(string(data))
	if media.Maker != "" {
		media.Title = strings.Replace(media.Title, media.Maker, "", 1)
	}
	media.Title = strings.Replace(media.Title, " ", "-", 10)
	media.M3u8Url = parseMadouM3U8Url(media.M3u8Url)
	return &media
}
