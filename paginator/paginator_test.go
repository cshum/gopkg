package paginator

import "testing"

func TestPaginator_GetTo(t *testing.T) {
	p := New(13)
	p.SetItemCount(20)
	p.Page = 1
	if p.GetTo() != 13 {
		t.Error("invalid")
	}
	p.Page = 2
	if p.GetTo() != 20 {
		t.Error("out of range")
	}
}
