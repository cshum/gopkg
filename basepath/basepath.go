package basepath

import (
	"github.com/kardianos/osext"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var once sync.Once
var basePath string

const stub = "!#"

// Init init base path from relative caller dir
func Init(elem ...string) {
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
			skip := 3
			if len(elem) == 1 && elem[0] == stub {
				elem[0] = "./"
				skip = 4
			}
			_, callerFileName, _, _ := runtime.Caller(skip)
			callerDir, err := filepath.Abs(filepath.Dir(callerFileName))
			if err != nil {
				panic(err)
			}
			basePath = filepath.Join(append([]string{callerDir}, elem...)...)
		}
	})
}

// Get abs project base path
func Get() string {
	Init(stub)
	return basePath
}

// Resolve get abs path from project base
func Resolve(path string) string {
	Init(stub)
	return filepath.Join(basePath, path)
}
