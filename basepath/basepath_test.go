package basepath

import (
	"fmt"
	"sync"
	"testing"
)

func TestBasePath(t *testing.T) {
	fmt.Println(Get())
	fmt.Println(Resolve("./test"))
	once = sync.Once{}
	Init("../")
	fmt.Println(Get())
	fmt.Println(Resolve("./test"))
}
