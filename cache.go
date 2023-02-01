package alfred

import (
	"encoding/json"
	"github.com/GitbookIO/diskache"
	"log"
)

var cache *diskache.Diskache

func CacheData[T any](k string, loader func() T) (T, bool) {
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
		err = cache.Set(k, b)
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
