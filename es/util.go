package es

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/olivere/elastic"
)

func GetIndicesByAlias(ctx context.Context, es *elastic.Client, name string) ([]string, error) {
	aliases, err := es.Aliases().Do(ctx)
	if err != nil {
		return []string{}, err
	}
	indices := aliases.IndicesByAlias(name)
	if indices == nil {
		return []string{}, nil
	}
	return indices, nil
}

func UpdateAlias(ctx context.Context, es *elastic.Client, name string, index string) error {
	indices, err := GetIndicesByAlias(ctx, es, name)
	if err != nil {
		return err
	}
	actions := []elastic.AliasAction{elastic.NewAliasAddAction(name).Index(index)}
	for _, idx := range indices {
		// delete alias if exists
		actions = append(actions, elastic.NewAliasRemoveAction(name).Index(idx))
	}
	if _, err := es.Alias().Action(actions...).Do(ctx); err != nil {
		return err
	}
	return nil
}

func CleanupOldIndices(ctx context.Context, es *elastic.Client, prefix, except string, keep int) error {
	indices, err := es.
		CatIndices().
		Index(prefix + "*").
		Sort("index:desc").
		Do(ctx)
	if err != nil {
		return err
	}
	for i, idx := range indices {
		if keep > 0 && i >= keep && idx.Index != except {
			if _, err := es.DeleteIndex(idx.Index).Do(ctx); err != nil {
				return err
			}
		}
	}
	return nil
}

func WaitForGreenStatus(
	ctx context.Context, es *elastic.Client,
	interval time.Duration, indices ...string,
) error {
	if len(indices) == 0 {
		return nil
	}
	var done bool
	for !done {
		done = true
		result, err := es.
			CatIndices().
			Index(strings.Join(indices, ",")).
			Do(ctx)
		if err != nil {
			return err
		}
		for _, idx := range result {
			if idx.Health == "red" {
				return errors.New(idx.Index + " red status, non-recoverable")
			}
			if idx.Health != "green" {
				time.Sleep(interval)
				done = false
				continue
			}
		}
	}
	return nil
}
