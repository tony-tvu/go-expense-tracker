package config

import (
	"log"
	"time"

	"github.com/allegro/bigcache/v3"
	"gorm.io/gorm"
)

func CreateCache(db *gorm.DB) *bigcache.BigCache {
	cacheConfig := bigcache.Config{
		Shards:     16,
		LifeWindow: 10 * time.Second,
		// 0 = never clean expired caches
		CleanWindow:        0,
		MaxEntriesInWindow: 1000 * 10 * 60,
		MaxEntrySize:       300,
		Verbose:            true,
		HardMaxCacheSize:   0,
		OnRemove:           nil,
		OnRemoveWithReason: nil,
	}

	cache, err := bigcache.NewBigCache(cacheConfig)
	if err != nil {
		log.Fatal(err)
	}

	var exists bool
	db.Raw("SELECT EXISTS(SELECT 1 FROM configs) AS found").Scan(&exists)
	if !exists {
		db.Create(&Config{})
		return cache
	}

	// TODO: populate cache from existing configs in db
	cache.Set("my-unique-key", []byte("hello there"))

	return cache
}
