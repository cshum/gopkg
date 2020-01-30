package es

import (
	"context"
	"encoding/json"
	"sort"
	"strings"
	"time"

	"github.com/cshum/gopkg/strof"
	"github.com/cshum/gopkg/tinycache"
	"github.com/olivere/elastic"
	"go.uber.org/zap"
)

type CachedRequest struct {
	Elastic    *elastic.Client
	Cache      tinycache.Cache
	Prefix     string
	Key        string
	Threshold  time.Duration
	Refresh    time.Duration
	Expiration time.Duration
	Logger     *zap.Logger
}

type CachedPayload struct {
	Timestamp int64                 `json:"ts"`
	Result    *elastic.SearchResult `json:"res"`
}

func (r *CachedRequest) Do(
	ctx context.Context,
	indices []string,
	source *elastic.SearchSource,
) (*elastic.SearchResult, error) {
	if r.Key == "" {
		key, err := SearchCacheKey(indices, source)
		if err != nil {
			return nil, err
		}
		r.Key = r.Prefix + key
	}
	if cached, ts := r.getSearchCache(r.Key); cached != nil {
		elasped := time.Millisecond * time.Duration(Timestamp()-ts)
		if r.Refresh > 0 && elasped >= r.Refresh {
			go func() {
				if result, err := r.Elastic.Search(indices...).
					SearchSource(source).Do(context.Background()); err == nil {
					r.setSearchCache(r.Key, result)
				}
			}()
		}
		return cached, nil
	}
	result, err := r.Elastic.Search(indices...).SearchSource(source).Do(ctx)
	if err == nil && result.TookInMillis > int64(r.Threshold/time.Millisecond) {
		go r.setSearchCache(r.Key, result)
	}
	return result, err
}

func (r *CachedRequest) setSearchCache(key string, result *elastic.SearchResult) {
	if val, err := json.Marshal(&CachedPayload{
		Timestamp(), result,
	}); err == nil {
		if err := r.Cache.Set(key, val, r.Expiration); err != nil && r.Logger != nil {
			r.Logger.Error("es-cache", zap.Error(err))
		}
	}
	if r.Logger != nil {
		r.Logger.Debug("es-cache",
			zap.String("action", "setSearchCache"),
			zap.String("key", key))
	}
}

func (r *CachedRequest) getSearchCache(key string) (*elastic.SearchResult, int64) {
	if val, err := r.Cache.Get(key); err == nil && len(val) > 0 {
		cached := &CachedPayload{}
		if err := json.Unmarshal(val, cached); err == nil {
			return cached.Result, cached.Timestamp
		} else if err == tinycache.NotFound && r.Logger != nil {
			r.Logger.Error("es-cache", zap.Error(err))
		}
	}
	return nil, 0
}

func SearchCacheKey(indices []string, source *elastic.SearchSource) (string, error) {
	sort.Strings(indices)
	src, err := source.Source()
	if err != nil {
		return "", err
	}
	hash, err := strof.Hash(src)
	if err != nil {
		return "", err
	}
	key := "!es!" + strings.Join(indices, ",") + "!" + hash
	return key, nil
}
