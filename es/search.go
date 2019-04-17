package es

import (
	"context"
	"errors"
	"github.com/cshum/gopkg/paginator"
	"github.com/olivere/elastic"
)

type Middleware func(s *Search)
type QueryHandler func(ctx context.Context, q *elastic.BoolQuery) error
type FunctionScoreHandler func(ctx context.Context, q *elastic.FunctionScoreQuery) error
type SourceHandler func(ctx context.Context, s *elastic.SearchSource) error
type ResultHandler func(ctx context.Context, result *elastic.SearchResult) error

type Search struct {
	Client   *elastic.Client
	indices  []string
	queries  []QueryHandler
	fnscores []FunctionScoreHandler
	sources  []SourceHandler
	sorters  []elastic.Sorter
	results  []ResultHandler
}

func NewSearch(es *elastic.Client, indices ...string) *Search {
	return &Search{Client: es, indices: indices}
}

func (q *Search) Index(indices ...string) *Search {
	q.indices = append(q.indices, indices...)
	return q
}

func (q *Search) HandleQuery(fn QueryHandler) *Search {
	q.queries = append(q.queries, fn)
	return q
}

func (q *Search) Sort(field string, ascending bool) *Search {
	q.sorters = append(q.sorters, elastic.SortInfo{
		Field:     field,
		Ascending: ascending,
	})
	return q
}

func (q *Search) HandleFunctionScore(fn FunctionScoreHandler) *Search {
	q.fnscores = append(q.fnscores, fn)
	return q
}

func (q *Search) HandleSource(fn SourceHandler) *Search {
	q.sources = append(q.sources, fn)
	return q
}

func (q *Search) New(indices ...string) *Search {
	return NewSearch(q.Client, indices...)
}

func (q *Search) HandleResult(fn ResultHandler) *Search {
	q.results = append(q.results, fn)
	return q
}

func (q *Search) Use(fn Middleware) *Search {
	fn(q)
	return q
}

func (q *Search) DoSource(
	ctx context.Context, p *paginator.Paginator,
) (*elastic.SearchSource, error) {
	bq := elastic.NewBoolQuery()
	if q.queries != nil {
		for _, query := range q.queries {
			if err := query(ctx, bq); err != nil {
				return nil, err
			}
		}
	}
	s := elastic.NewSearchSource()
	if p != nil && len(q.fnscores) > 0 {
		fsq := elastic.NewFunctionScoreQuery().
			BoostMode("replace").
			Query(bq)
		for _, fnscore := range q.fnscores {
			if err := fnscore(ctx, fsq); err != nil {
				return nil, err
			}
		}
		s.Query(fsq)
	} else {
		s.Query(bq)
	}
	if len(q.sorters) > 0 {
		s.SortBy(q.sorters...)
	}
	if p != nil {
		if p.GetOffset()+p.GetLimit() >= 10000 {
			return nil, errors.New("page exceeded maximum")
		}
		s.From(p.GetOffset())
		s.Size(p.GetLimit())
	} else {
		s.Size(0)
	}
	if q.sources != nil {
		for _, search := range q.sources {
			if err := search(ctx, s); err != nil {
				return nil, err
			}
		}
	}
	return s, nil
}

func (q *Search) Do(
	ctx context.Context, p *paginator.Paginator,
) (*elastic.SearchResult, error) {
	ss, err := q.DoSource(ctx, p)
	if err != nil {
		return nil, err
	}
	result, err := q.Client.
		Search(q.indices...).
		SearchSource(ss).
		Do(ctx)
	if err != nil {
		return result, err
	}
	if p != nil {
		p.SetCount(result.TotalHits())
	}
	if q.results != nil {
		for _, res := range q.results {
			if err := res(ctx, result); err != nil {
				return result, err
			}
		}
	}
	return result, nil
}
