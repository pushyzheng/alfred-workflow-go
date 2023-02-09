package alfred

import (
	"encoding/json"
	"github.com/pushyzheng/diskache"
	"github.com/sirupsen/logrus"
	"log"
)

var (
	cache    *diskache.Diskache
	cacheLog = NewFileLogger("cache")
)

func CacheData[T any](k string, loader func() T) (T, bool) {
	return CacheExpiredData(k, -1, loader)
}

func CacheExpiredData[T any](k string, expired int64, loader func() T) (T, bool) {
	cacheLog.WithFields(logrus.Fields{
		"key":     k,
		"expired": expired,
	}).Info("CacheExpiredData")
	data, exists := cache.Get(k)
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
	cacheLog.WithFields(logrus.Fields{
		"key":     k,
		"expired": expired,
		"loaded":  !exists,
	}).Info("CacheExpiredData succeed")
	return res, exists
}

func GetCacheData(k string) ([]byte, bool) {
	cacheLog.WithFields(logrus.Fields{
		"key": k,
	}).Info("GetCacheData")
	return cache.Get(k)
}

func SetCacheJsonData(k string, expired int64, data any) {
	cacheLog.WithFields(logrus.Fields{
		"key":     k,
		"expired": expired,
	}).Info("SetCacheJsonData")
	b, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	SetCacheData(k, expired, b)
}

func SetCacheData(k string, expired int64, data []byte) {
	var err error
	if expired == -1 {
		err = cache.Set(k, data)
	} else {
		err = cache.SetExpired(k, data, expired)
	}
	if err != nil {
		log.Printf("error: set cache error, key = %s", k)
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
