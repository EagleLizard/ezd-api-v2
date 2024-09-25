package constants

import (
	"log"
	"path/filepath"
	"runtime"
)

/*
	How many dirs the root dir is
		above the current file
*/

const parentLvls = 4

func BaseDir() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Error getting BaseDir")
	}
	baseDir := filename
	for i := 0; i < parentLvls; i++ {
		baseDir = filepath.Dir(baseDir)
	}
	return baseDir
}
