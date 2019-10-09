package paginator

import (
	"errors"
	"math"
)

// Paginator for embedding pagination
type Paginator struct {
	Page      int `json:"page" validate:"min=1"`
	Size      int `json:"size" validate:"min=1,max=100"`
	itemcount int
	pagecount int
}

func New(size int) *Paginator {
	if size < 1 {
		panic(errors.New("size must be non-zero positive integer"))
	}
	return &Paginator{Page: 1, Size: size}
}

func (p *Paginator) GetOffset() int {
	return (p.Page - 1) * p.Size
}

func (p *Paginator) GetFrom() int {
	from := p.GetOffset()
	if cnt := p.itemcount; from > cnt {
		return cnt
	}
	return from
}

func (p *Paginator) GetTo() int {
	to := p.GetOffset() + p.Size
	if cnt := p.itemcount; to > cnt {
		return cnt
	}
	return to
}

func (p *Paginator) GetLimit() int {
	return p.Size
}

func (p *Paginator) GetItemCount() int {
	return p.itemcount
}

func (p *Paginator) GetPageCount() int {
	return p.pagecount
}

func (p *Paginator) HasNext() bool {
	return p.pagecount >= p.Page+1
}

func (p *Paginator) Next() bool {
	if !p.HasNext() {
		return false
	}
	p.Page++
	return true
}

func (p *Paginator) SetItemCount(count int) {
	p.itemcount = count
	p.pagecount = int(math.Ceil(float64(p.itemcount) / float64(p.Size)))
}

// Pagination for json response
type Pagination struct {
	HasNext   bool `json:"has_next"`
	ItemCount int  `json:"item_count"`
	PageCount int  `json:"page_count"`
	Page      int  `json:"page"`
	Size      int  `json:"size"`
}

// Pagination create struct from pagecount results itemcount
func (p *Paginator) Pagination() *Pagination {
	return &Pagination{
		HasNext:   p.HasNext(),
		ItemCount: p.itemcount,
		PageCount: p.pagecount,
		Page:      p.Page,
		Size:      p.Size,
	}
}
