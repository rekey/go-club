package mover

import (
	"encoding/xml"
	"os"
	"path/filepath"

	"github.com/rekey/go-club/parser"
)

type Actor struct {
	Name string `xml:"name,omitempty"`
	Type string `xml:"type,omitempty"`
}

type Movie struct {
	XMLName       xml.Name `xml:"movie,omitempty"`
	Title         string   `xml:"title,omitempty"`
	Originaltitle string   `xml:"originaltitle,omitempty"`
	Genre         []string `xml:"genre,omitempty"`
	Studio        string   `xml:"studio,omitempty"`
	Tag           []string `xml:"tag,omitempty"`
	Poster        string   `xml:"poster,omitempty"`
	Thumb         string   `xml:"thumb,omitempty"`
	Fanart        string   `xml:"fanart,omitempty"`
	Maker         string   `xml:"maker,omitempty"`
	Cover         string   `xml:"cover,omitempty"`
	Actor         []Actor  `xml:"actor,omitempty"`
}

func ParseMediaToFile(media *parser.Media, dir string, thumbPath string) error {
	nfoName := media.Title + ".nfo"
	nfoPath := filepath.Join(dir, nfoName)
	movie := &Movie{
		Title:         media.Title,
		Originaltitle: media.Title,
		Tag:           media.Tags,
		Studio:        media.Maker,
		Maker:         media.Maker,
		Genre:         media.Tags,
	}
	if media.Actor != "" {
		movie.Actor = []Actor{
			{
				Type: "Actor",
				Name: media.Actor,
			},
		}
	}
	if thumbPath != "" {
		if _, err := os.Stat(thumbPath); err == nil {
			movie.Thumb = filepath.Base(thumbPath)
			movie.Cover = filepath.Base(thumbPath)
			movie.Poster = filepath.Base(thumbPath)
		}
	}
	data, err := xml.MarshalIndent(movie, "", "  ")
	if err != nil {
		return err
	}
	data = append([]byte(xml.Header), data...)
	err = os.WriteFile(nfoPath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}
