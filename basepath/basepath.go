package basepath

import (
	"github.com/cshum/gopkg/util"
	"github.com/kardianos/osext"
	"path/filepath"
	"strings"
)

type BasePath struct {
	path string
}

// New init base path
func New(rel string) *BasePath {
	exeDir, err := osext.ExecutableFolder()
	if err != nil {
		panic(err)
	}
	path, err := filepath.Abs(exeDir)
	if err != nil {
		panic(err)
	}
	if strings.HasSuffix(path, "/exe") ||
		strings.HasSuffix(path, "/T") {
		// execute via go run or go test
		callerDir, err := util.CallerDir()
		if err != nil {
			panic(err)
		}
		path = filepath.Join(callerDir, rel)
	}
	return &BasePath{path}
}

// Get abs project base path
func (b *BasePath) Get() string {
	return b.path
}

// Resolve get abs path from project base
func (b *BasePath) Resolve(path string) string {
	return filepath.Join(b.path, path)
}
