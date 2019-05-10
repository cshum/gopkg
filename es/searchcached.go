package es

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/cshum/gopkg/util"
	"github.com/go-redis/redis"
	"github.com/olivere/elastic"
	"go.uber.org/zap"
	"sort"
	"strings"
	"time"
)

type CachedRequest struct {
	Elastic    *elastic.Client
	Redis      *redis.Client
	Key        string
	Threshold  time.Duration
	Refresh    time.Duration
	Expiration time.Duration
	Logger     *zap.Logger
}

type SearchCache struct {
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
		r.Key = key
	}
	if cached, ts := r.getSearchCache(r.Key); cached != nil {
		elasped := time.Millisecond * time.Duration(util.Timestamp()-ts)
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
	if val, err := json.Marshal(&SearchCache{
		util.Timestamp(), result,
	}); err == nil {
		if _, err := r.Redis.Set(key, val, r.Expiration).Result(); err != nil && r.Logger != nil {
			r.Logger.Error("redis", zap.Error(err))
		}
	}
	if r.Logger != nil {
		r.Logger.Debug("redis",
			zap.String("action", "setSearchCache"),
			zap.String("key", key))
	}
}

func (r *CachedRequest) getSearchCache(key string) (*elastic.SearchResult, int64) {
	if val, err := r.Redis.Get(key).Result(); err == nil {
		cached := &SearchCache{}
		if err := json.Unmarshal([]byte(val), cached); err == nil {
			return cached.Result, cached.Timestamp
		} else if err != redis.Nil && r.Logger != nil {
			r.Logger.Error("redis", zap.Error(err))
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
	bytes, err := json.Marshal(src)
	if err != nil {
		return "", err
	}
	hash := md5.Sum(bytes)
	key := "!es!" + strings.Join(indices, ",") + "!" + hex.EncodeToString(hash[:])
	return key, nil
}
