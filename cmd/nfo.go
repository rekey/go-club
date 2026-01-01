package main

import (
	"os"
	"path/filepath"

	"github.com/rekey/go-club/mover"
	"github.com/rekey/go-club/parser"
)

func main() {
	cwd, _ := os.Getwd()
	mover.ParseMediaToFile(&parser.Media{
		Title: "test",
		Tags:  []string{"a", "b"},
		Actor: "c",
	}, cwd, filepath.Join(cwd, "poster.jpg"))
}
