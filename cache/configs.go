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

var PAGE_LIMIT int64 = 500

type ConfigsInput struct {
	RegistrationEnabled bool `json:"registration_enabled"`
	QuotaEnabled        bool `json:"quota_enabled"`
	QuotaLimit          int  `json:"quota_limit"`
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
			{Key: "quota_enabled", Value: true},
			{Key: "quota_limit", Value: 10},
			{Key: "registration_enabled", Value: false},
			{Key: "page_limit", Value: PAGE_LIMIT},
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

func (c *Configs) UpdateConfigsCache(ctx context.Context, db *database.MongoDb, input *ConfigsInput) error {
	var configs *models.Config
	if err := db.Configs.FindOne(ctx, bson.M{}).Decode(&configs); err != nil {
		return err
	}
	_, err := db.Configs.UpdateOne(
		ctx,
		bson.M{"_id": configs.ID},
		bson.M{
			"$set": bson.M{
				"quota_enabled":        input.QuotaEnabled,
				"quota_limit":          input.QuotaLimit,
				"registration_enabled": input.RegistrationEnabled,
				"updated_at":           time.Now(),
			}},
	)
	if err != nil {
		return err
	}

	var updatedConfigs *models.Config
	if err := db.Configs.FindOne(ctx, bson.M{}).Decode(&updatedConfigs); err != nil {
		return err
	}

	b, err := json.Marshal(updatedConfigs)
	if err != nil {
		return err
	}

	err = c.Cache.Set("configs", b)
	if err != nil {
		return err
	}
	return nil
}
