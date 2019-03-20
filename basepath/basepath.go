package basepath

import (
	"github.com/cshum/gopkg/util"
	"github.com/kardianos/osext"
	"path/filepath"
	"strings"
	"sync"
)

var once sync.Once
var basePath string

// Init init base path from relative caller dir
func Init(rel string) {
	once.Do(func() {
		exeDir, err := osext.ExecutableFolder()
		if err != nil {
			panic(err)
		}
		basePath, err = filepath.Abs(exeDir)
		if err != nil {
			panic(err)
		}
		if strings.HasSuffix(basePath, "/exe") ||
			strings.HasSuffix(basePath, "/T") {
			// execute via go run or go test
			callerDir, err := util.CallerDir()
			if err != nil {
				panic(err)
			}
			basePath = filepath.Join(callerDir, rel)
		}
	})
}

// Get abs project base path
func Get() string {
	Init("./")
	return basePath
}

// Resolve get abs path from project base
func Resolve(path string) string {
	Init("./")
	return filepath.Join(basePath, path)
}
