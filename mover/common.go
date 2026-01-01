package mover

import (
	"github.com/rekey/go-club/common"
	"github.com/rekey/go-club/env"
)

func init() {
	DownloadDir := env.DownloadDir
	common.CreateDir(DownloadDir)
	common.CreateDir(env.DownloadResultsDir)
}
