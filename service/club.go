package service

import (
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/canhlinh/hlsdl"
)

func parseClub(url string) *Media {
	res, err := http.Get(url)
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
	var media = Media{}
	media.Title = doc.Find(".article-title").First().Text()
	media.Maker = doc.Find(".article-meta .item-3").First().Text()
	media.Tags = doc.Find(".article-tags a").Map(func(i int, s *goquery.Selection) string {
		return s.Text()
	})
	media.Url, _ = doc.Find(".article-content iframe").First().Attr("src")
	// data, _ := json.Marshal(media)
	// log.Println(string(data))
	if media.Maker != "" {
		media.Title = strings.Replace(media.Title, media.Maker, "", 1)
	}
	media.Title = strings.Replace(media.Title, " ", "-", 10)
	filename := media.Title + ".mp4"
	media.Url = down(PTClub, media.Url, filename)
	return &media
}

func downClub(url string, filename string) string {
	res, err := http.Get(url)
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
	m3u8Url := "https://dash.madou.club" + m3u8 + "?token=" + token
	log.Println("start download")
	hlsDL := hlsdl.New(m3u8Url, nil, "download", filename, 20, true)
	file, err := hlsDL.Download()
	if err != nil {
		log.Println(err)
		return ""
	}
	log.Println(file, "done")
	return ""
}
