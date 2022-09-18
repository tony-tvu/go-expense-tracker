package cache

import (
	"encoding/json"
	"log"

	"github.com/allegro/bigcache"
	"github.com/tony-tvu/goexpense/entity"
	"gorm.io/gorm"
)

// These configuration values are accessed regularly in numerous locations.
// To avoid making excessive calls to the database, the configs cache is
// populated on startup using the existing configs from the db and the cached
// object can be accessed wherever needed. Whenever a config value is updated
// within the db, the cached configs will also get updated.
type Configs struct {
	Cache *bigcache.BigCache
}

func (c *Configs) InitConfigsCache(db *gorm.DB) {
	cacheConfig := bigcache.Config{
		Shards:     16,
		LifeWindow: 0,
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
		db.Create(&entity.Config{})
	}

	var conf *entity.Config
	if result := db.Exec("SELECT * FROM configs LIMIT 1;").First(&conf); result.Error != nil {
		log.Fatal(err)
	}

	b, err := json.Marshal(conf)
	if err != nil {
		log.Fatal(err)
	}

	cache.Set("configs", b)
	c.Cache = cache
}

func (c *Configs) GetConfigs() (*entity.Config, error) {
	var conf entity.Config
	b, err := c.Cache.Get("configs")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
