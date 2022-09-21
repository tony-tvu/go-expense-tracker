package cache

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/allegro/bigcache"
	"github.com/tony-tvu/goexpense/database"
	"github.com/tony-tvu/goexpense/models"
	"go.mongodb.org/mongo-driver/bson"
)

// These configuration values are accessed regularly in numerous locations.
// To avoid making excessive calls to the database, the configs cache is
// populated on startup using the existing configs from the db and the cached
// object can be accessed wherever needed. Whenever a config value is updated
// within the db, the cached configs will also get updated.
type Configs struct {
	Cache *bigcache.BigCache
}

func (c *Configs) InitConfigsCache(ctx context.Context, db *database.MongoDb) {
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

	count, err := db.Configs.CountDocuments(ctx, bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		// init default configs
		doc := &bson.D{
			{Key: "access_token_exp", Value: 900},
			{Key: "refresh_token_exp", Value: 3600},
			{Key: "quota_enabled", Value: true},
			{Key: "quota_limit", Value: 10},
			{Key: "tasks_enabled", Value: true},
			{Key: "tasks_interval", Value: 60},
			{Key: "registration_enabled", Value: false},
			{Key: "created_at", Value: time.Now()},
			{Key: "updated_at", Value: time.Now()},
		}
		if _, err := db.Configs.InsertOne(ctx, doc); err != nil {
			log.Fatal(err)
		}
	}

	var configs *models.Config
	if err = db.Configs.FindOne(ctx, bson.M{}).Decode(&configs); err != nil {
		log.Fatal(err)
	}

	b, err := json.Marshal(configs)
	if err != nil {
		log.Fatal(err)
	}

	cache.Set("configs", b)
	c.Cache = cache
}

func (c *Configs) GetConfigs() (*models.Config, error) {
	var configs models.Config
	b, err := c.Cache.Get("configs")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &configs)
	if err != nil {
		return nil, err
	}

	return &configs, nil
}
