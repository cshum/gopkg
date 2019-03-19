package basepath

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	fmt.Println(New("../").Get())
	fmt.Println(New("").Get())
}
