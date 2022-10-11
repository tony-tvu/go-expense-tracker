package cache

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/allegro/bigcache"
	"github.com/joho/godotenv"
	"github.com/tony-tvu/goexpense/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// These configuration values are accessed regularly in numerous locations.
// To avoid making excessive calls to the database, the configs cache is
// populated on startup using the existing configs from the db and the cached
// object can be accessed wherever needed. Whenever a config value is updated
// within the db, the cached configs will also get updated.
type ConfigsCache struct {
	Cache *bigcache.BigCache
}

type Config struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`

	// If false, users will not be able to create accounts from UI or handler routes
	RegistrationEnabled bool   `json:"registration_enabled" bson:"registration_enabled"`
	PageLimit           int64  `json:"page_limit" bson:"page_limit"`
	TellerApplicationID string `json:"teller_application_id" bson:"teller_application_id"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}


var PAGE_LIMIT int64 = 500

type ConfigsInput struct {
	RegistrationEnabled bool `json:"registration_enabled"`
}

func (c *ConfigsCache) InitConfigsCache(ctx context.Context, db *db.MongoDb) {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("no .env file found")
	}

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
			{Key: "registration_enabled", Value: false},
			{Key: "page_limit", Value: PAGE_LIMIT},
			{Key: "teller_application_id", Value: os.Getenv("TELLER_APPLICATION_ID")},
			{Key: "created_at", Value: time.Now()},
			{Key: "updated_at", Value: time.Now()},
		}
		if _, err := db.Configs.InsertOne(ctx, doc); err != nil {
			log.Fatal(err)
		}
	}

	var configs *Config
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

func (c *ConfigsCache) GetConfigs() (*Config, error) {
	var config Config
	b, err := c.Cache.Get("configs")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *ConfigsCache) UpdateConfigsCache(ctx context.Context, db *db.MongoDb, input *ConfigsInput) error {
	var config *Config
	if err := db.Configs.FindOne(ctx, bson.M{}).Decode(&config); err != nil {
		return err
	}
	_, err := db.Configs.UpdateOne(
		ctx,
		bson.M{"_id": config.ID},
		bson.M{
			"$set": bson.M{
				"registration_enabled": input.RegistrationEnabled,
				"updated_at":           time.Now(),
			}},
	)
	if err != nil {
		return err
	}

	var updatedConfig *Config
	if err := db.Configs.FindOne(ctx, bson.M{}).Decode(&updatedConfig); err != nil {
		return err
	}

	b, err := json.Marshal(updatedConfig)
	if err != nil {
		return err
	}

	err = c.Cache.Set("configs", b)
	if err != nil {
		return err
	}
	return nil
}
