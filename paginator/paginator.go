package paginator

import (
	"errors"
	"math"
)

// Paginator for embedding pagination
type Paginator struct {
	Page      int `json:"page" validate:"min=1"`
	Size      int `json:"size" validate:"min=1,max=100"`
	itemcount int64
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
	return p.GetOffset()
}

func (p *Paginator) GetTo() int {
	return p.GetOffset() + p.Size
}

func (p *Paginator) GetLimit() int {
	return p.Size
}

func (p *Paginator) GetItemCount() int64 {
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

func (p *Paginator) SetItemCount(count int64) {
	p.itemcount = count
	p.pagecount = int(math.Ceil(float64(p.itemcount) / float64(p.Size)))
}

// Pagination for json response
type Pagination struct {
	HasNext   bool  `json:"has_next"`
	ItemCount int64 `json:"item_count"`
	PageCount int   `json:"page_count"`
	Page      int   `json:"page"`
	Size      int   `json:"size"`
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
