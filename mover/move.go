package mover

import (
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/rekey/go-club/common"
	"github.com/rekey/go-club/env"
	"github.com/rekey/go-club/parser"
)

func Move(media *parser.Media, dir string) error {
	tmpDir := filepath.Join(env.DownloadTmpDir, dir)
	resultDir := filepath.Join(env.DownloadResultsDir, dir)
	common.CreateDir(resultDir)
	log.Println("move", tmpDir, "to", resultDir)
	files, err := os.ReadDir(tmpDir)
	log.Println("read tmp dir error", err, files)
	if err != nil {
		return err
	}
	for _, file := range files {
		filepath.Join(tmpDir, file.Name())
		tmpFile := filepath.Join(tmpDir, file.Name())
		resultFile := filepath.Join(resultDir, file.Name())
		log.Println("move", tmpFile, "to", resultFile)
		os.Rename(tmpFile, resultFile)
	}
	posterFile := "poster" + path.Ext(media.Thumb)
	posterPath := filepath.Join(resultDir, posterFile)
	if media.Thumb != "" {
		downloadThumb(media, posterPath)
	}
	err = ParseMediaToFile(media, resultDir, posterPath)
	if err != nil {
		return err
	}
	return nil
}
