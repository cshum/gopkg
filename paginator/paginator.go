package paginator

import (
	"errors"
	"math"

	"github.com/volatiletech/sqlboiler/queries/qm"
)

// Paginator for embedding pagination
type Paginator struct {
	Page  int `json:"page" validate:"min=1"`
	Size  int `json:"per" validate:"min=1,max=100"`
	count int64
	total int
}

// New create Paginator of page size
func New(size int) *Paginator {
	if size < 1 {
		panic(errors.New("size must be non-zero positive integer"))
	}
	return &Paginator{Page: 1, Size: size}
}

// GetOffset calculate paginate offset
func (p *Paginator) GetOffset() int {
	return (p.Page - 1) * p.Size
}

// GetLimit calculate paginate limit
func (p *Paginator) GetLimit() int {
	return p.Size
}

// GetCount result count
func (p *Paginator) GetCount() int64 { //sqlboiler uses int64
	return p.count
}

// GetTotal total page
func (p *Paginator) GetTotal() int {
	return p.total
}

// HasNext has next page
func (p *Paginator) HasNext() bool {
	return p.total >= p.Page+1
}

// SetCount set result count
func (p *Paginator) SetCount(count int64) {
	p.count = count
	p.total = int(math.Ceil(float64(p.count) / float64(p.Size)))
}

// PaginatorMods for sqlboiler
func (p *Paginator) PaginatorMods(q ...qm.QueryMod) []qm.QueryMod {
	var query []qm.QueryMod
	if len(q) > 0 {
		query = append(query, q...)
	}
	query = append(
		query,
		qm.Offset(p.GetOffset()),
		qm.Limit(p.GetLimit()),
	)
	return query
}

// Pagination for pagination json response
type Pagination struct {
	HasNext     bool `json:"has_next"`
	Total       int  `json:"total_pages"`
	CurrentPage int  `json:"current_page"`
	Size        int  `json:"per"`
}

// Pagination create struct from total results count
func (p *Paginator) Pagination() *Pagination {
	return &Pagination{
		HasNext:     p.HasNext(),
		Total:       p.total,
		CurrentPage: p.Page,
		Size:        p.Size,
	}
}
