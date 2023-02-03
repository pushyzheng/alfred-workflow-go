package alfred

import (
	"encoding/json"
	"github.com/pushyzheng/diskache"
	"log"
)

var cache *diskache.Diskache

func CacheData[T any](k string, loader func() T) (T, bool) {
	return CacheExpiredData(k, -1, loader)
}

func CacheExpiredData[T any](k string, expired int64, loader func() T) (T, bool) {
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
		b, err := json.Marshal(res)
		if err != nil {
			panic(err)
		}
		if expired == -1 {
			err = cache.Set(k, b)
		} else {
			err = cache.SetExpired(k, b, expired)
		}
		if err != nil {
			log.Printf("error: set cache error, key = %s", k)
		}
	}
	return res, exists
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
