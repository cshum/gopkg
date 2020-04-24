package basepath

import (
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/kardianos/osext"
)

type OnFileLoadedHandler func(path string, isPreloaded bool)

var once sync.Once
var basePath string
var fileMap = map[string][]byte{}
var l = sync.RWMutex{}
var onFileLoadedHandlers []OnFileLoadedHandler

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
			strings.Contains(basePath, "/T") {
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

func OnFileLoaded(fn OnFileLoadedHandler) {
	l.Lock()
	defer l.Unlock()
	onFileLoadedHandlers = append(onFileLoadedHandlers, fn)
}

func LoadFile(path string) ([]byte, error) {
	l.RLock()
	defer l.RUnlock()
	abspath := Resolve(path)
	var result []byte
	preloaded := false
	if bytes, ok := fileMap[abspath]; ok {
		result = bytes
		preloaded = true
	} else {
		bytes, err := ioutil.ReadFile(abspath)
		if err != nil {
			return bytes, err
		}
		result = bytes
	}
	if onFileLoadedHandlers != nil {
		for _, fn := range onFileLoadedHandlers {
			fn(abspath, preloaded)
		}
	}
	return result, nil
}

func LoadFileString(path string) (string, error) {
	bytes, err := LoadFile(path)
	return string(bytes), err
}

func PreloadFile(path string, data []byte) {
	l.Lock()
	defer l.Unlock()
	fileMap[Resolve(path)] = data
}
