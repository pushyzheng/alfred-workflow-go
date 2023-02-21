package alfred

import (
	"encoding/json"
	"log"

	"github.com/pushyzheng/diskache"
	"github.com/sirupsen/logrus"
)

var (
	cache    *diskache.Diskache
	cacheLog = NewFileLogger("cache")
)

// CacheData Set durable cache data
func CacheData[T any](k string, loader func() T) (T, bool) {
	return CacheExpiredData(k, -1, loader)
}

// CacheExpiredData Get data from local cache, and provide a loader method,
// when the cache does not exist, it will call the loader method to get the latest data
// and rewrite it to the local cache.
// when expired is set to -1, it means the cache won't expire.
func CacheExpiredData[T any](k string, expired int64, loader func() T) (T, bool) {
	data, exists := cache.Get(k)
	defer func() {
		cacheLog.WithFields(logrus.Fields{
			"key":     k,
			"expired": expired,
			"loaded":  !exists,
		}).Info("Cache expired data succeed")
	}()
	var res T
	var err error

	if exists {
		err = json.Unmarshal(data, &res)
		if err != nil {
			panic(err)
		}
	} else {
		res = loader()
		SetCacheJsonData(k, expired, res)
	}
	return res, exists
}

func CacheDataRefresh[T any](k string, loader func() T) {
	cacheLog.WithFields(logrus.Fields{
		"key": k,
	}).Info("cache data refresh")
	res := loader()
	SetCacheJsonData(k, -1, res)
}

// GetCacheData Get data from local cache
func GetCacheData(k string) ([]byte, bool) {
	cacheLog.WithFields(logrus.Fields{
		"key": k,
	}).Info("get cache data")
	return cache.Get(k)
}

// SetCacheJsonData Set the data to the local cache,
// and serialize the structure at the same time automatically.
func SetCacheJsonData(k string, expired int64, data any) {
	cacheLog.WithFields(logrus.Fields{
		"key":     k,
		"expired": expired,
	}).Info("set cache json data")
	b, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	SetCacheData(k, expired, b)
}

// SetCacheData Set data to local cache
func SetCacheData(k string, expired int64, data []byte) {
	var err error
	if expired == -1 {
		err = cache.Set(k, data)
	} else {
		err = cache.SetExpired(k, data, expired)
	}
	if err != nil {
		cacheLog.WithField("key", k).Errorf("set cache error")
	}
}

func init() {
	opts := diskache.Opts{
		Directory: "cache",
	}
	var err error
	cache, err = diskache.New(&opts)
	if err != nil {
		log.Fatalln(err)
	}
}
