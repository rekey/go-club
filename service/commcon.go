package service

import "log"

type Media struct {
	Title string
	Maker string
	Tags  []string
	Url   string
}

type ParseType int

const (
	PTClub ParseType = iota
)

func (parseType ParseType) IsValid() bool {
	return parseType >= PTClub && parseType <= PTClub
}

func Parse(parseType ParseType, url string) {
	switch parseType {
	case PTClub:
		parseClub(url)
	default:
		log.Println("don't know", parseType)
	}
}

func down(parseType ParseType, url string, filename string) string {
	switch parseType {
	case PTClub:
		return downClub(url, filename)
	default:
		log.Println("don't know", parseType)
		return ""
	}
}
