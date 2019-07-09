package es

import (
	"context"
	"encoding/json"
	"github.com/cshum/gopkg/util"
	"github.com/go-redis/redis"
	"github.com/olivere/elastic"
	"go.uber.org/zap"
	"sort"
	"strings"
	"time"
)

type Cache interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte, ttl time.Duration) error
}

type RedisCache struct {
	Client *redis.Client
}

func (r *RedisCache) Get(key string) ([]byte, error) {
	res, err := r.Client.Get(key).Result()
	return []byte(res), err
}

func (r *RedisCache) Set(key string, value []byte, ttl time.Duration) error {
	_, err := r.Client.Set(key, value, ttl).Result()
	return err
}

type CachedRequest struct {
	Elastic    *elastic.Client
	Cache      Cache
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
	if val, err := json.Marshal(&CachedPayload{
		util.Timestamp(), result,
	}); err == nil {
		if err := r.Cache.Set(key, val, r.Expiration); err != nil && r.Logger != nil {
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
	if val, err := r.Cache.Get(key); err == nil && len(val) > 0 {
		cached := &CachedPayload{}
		if err := json.Unmarshal(val, cached); err == nil {
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
	hash, err := util.ToHash(src)
	if err != nil {
		return "", err
	}
	key := "!es!" + strings.Join(indices, ",") + "!" + hash
	return key, nil
}
