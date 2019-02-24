package unique

import (
	"testing"
)

func TestInts(t *testing.T) {
	if ints := Ints([]int{0, 0, 1, 2}); len(ints) != 3 {
		t.Fatalf("1: Unexpected length on the returned slice: %v", len(ints))
	} else if ints[0] != 0 || ints[1] != 1 || ints[2] != 2 {
		t.Fatalf("1: Unexpected slice contents: %v", ints)
	}

	if ints := Ints([]int{0, 1, 2}); len(ints) != 3 {
		t.Fatalf("2: Unexpected length on the returned slice: %v", len(ints))
	} else if ints[0] != 0 || ints[1] != 1 || ints[2] != 2 {
		t.Fatalf("2: Unexpected slice contents: %v", ints)
	}

	if ints := Ints([]int{}); len(ints) != 0 {
		t.Fatalf("3: Unexpected length on the returned slice: %v", len(ints))
	}
}

func TestStrings(t *testing.T) {
	if strings := Strings([]string{"a", "a", "b", "c"}); len(strings) != 3 {
		t.Fatalf("1: Unexpected length on the returned slice: %v", len(strings))
	} else if strings[0] != "a" || strings[1] != "b" || strings[2] != "c" {
		t.Fatalf("1: Unexpected slice contents: %v", strings)
	}

	if strings := Strings([]string{"a", "b", "c"}); len(strings) != 3 {
		t.Fatalf("2: Unexpected length on the returned slice: %v", len(strings))
	} else if strings[0] != "a" || strings[1] != "b" || strings[2] != "c" {
		t.Fatalf("2: Unexpected slice contents: %v", strings)
	}

	if strings := Strings([]string{}); len(strings) != 0 {
		t.Fatalf("3: Unexpected length on the returned slice: %v", len(strings))
	}
}
