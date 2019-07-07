package es

import (
	"context"
	"encoding/json"

	"github.com/olivere/elastic"
)

type Index struct {
	Elastic *elastic.Client
	Index   string
	Body    IndexBody
}

type IndexBody struct {
	Settings *json.RawMessage            `json:"settings"`
	Mappings map[string]*json.RawMessage `json:"mappings"`
}

func NewIndex(es *elastic.Client, name string) *Index {
	return &Index{
		Elastic: es,
		Index:   name,
		Body: IndexBody{
			Mappings: map[string]*json.RawMessage{},
		},
	}
}

func (idx *Index) SetSettings(settings string) {
	msg := json.RawMessage(settings)
	idx.Body.Settings = &msg
}

func (idx *Index) SetMapping(name, mapping string) {
	msg := json.RawMessage(mapping)
	idx.Body.Mappings[name] = &msg
}

func (idx *Index) ToJSONString() (string, error) {
	bytes, err := json.Marshal(idx.Body)
	return string(bytes), err
}

func (idx *Index) Create(ctx context.Context) error {
	exists, err := idx.Elastic.IndexExists(idx.Index).Do(ctx)
	if err != nil {
		return err
	}
	if exists {
		_, err := idx.Elastic.DeleteIndex(idx.Index).Do(ctx)
		if err != nil {
			return err
		}
	}
	body, err := idx.ToJSONString()
	if err != nil {
		return err
	}
	_, err = idx.Elastic.CreateIndex(idx.Index).Body(body).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}
