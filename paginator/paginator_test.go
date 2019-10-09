package paginator

import "testing"

func TestPaginator_GetFrom_GetTo(t *testing.T) {
	p := New(13)
	p.SetItemCount(20)
	p.Page = 1
	if p.GetTo() != 13 {
		t.Error("invalid")
	}
	p.Page = 2
	if p.GetFrom() != 13 {
		t.Error("invalid")
	}
	if p.GetTo() != 20 {
		t.Error("out of bound")
	}
	p.Page = 3
	if p.GetFrom() != 20 {
		t.Error("out of bound")
	}
	if p.GetTo() != 20 {
		t.Error("out of bound")
	}
}
