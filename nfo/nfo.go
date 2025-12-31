package nfo

import (
	"github.com/rekey/go-club/parser"

	"github.com/simulot/aspiratv/metadata/nfo"
)

func MediaToNFO(media *parser.Media, nfoPath string) {
	movie := &nfo.Movie{
		MediaInfo: nfo.MediaInfo{
			Title: media.Title,
			Tag:   media.Tags,
		},
	}
	movie.WriteNFO(nfoPath)
}
