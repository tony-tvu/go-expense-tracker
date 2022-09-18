package cache

import (
	"fmt"
	"log"
	"strconv"

	"github.com/allegro/bigcache"
	"github.com/tony-tvu/goexpense/entity"
	"gorm.io/gorm"
)

type ConfigCache struct {
	Cache *bigcache.BigCache
}

func (c *ConfigCache) InitConfigCache(db *gorm.DB) {
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

	cache.Set("registration_allowed", []byte(fmt.Sprint(conf.RegistrationAllowed)))
	cache.Set("quota_enforced", []byte(fmt.Sprint(conf.QuotaEnforced)))
	cache.Set("quota_limit", []byte(fmt.Sprint(conf.QuotaLimit)))

	c.Cache = cache
}

func (c *ConfigCache) GetRegistrationAllowed() (*bool, error) {
	value, err := c.Cache.Get("registration_allowed")
	if err != nil {
		return nil, err
	}

	valuebool, err := strconv.ParseBool(string(value))
	if err != nil {
		return nil, err
	}

	return &valuebool, nil
}
