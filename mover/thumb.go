package mover

import (
	"io"
	"log"
	"os"

	"github.com/rekey/go-club/common"
	"github.com/rekey/go-club/parser"
)

func downloadThumb(media *parser.Media, output string) error {
	log.Println("output", output)
	res, err := common.HttpGet(media.Thumb, common.GetUrlRefer(media.Url))
	log.Println("thumb download error", err)
	if err != nil {
		log.Println(err)
		return err
	}
	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return os.WriteFile(output, bodyBytes, 0644)
}
