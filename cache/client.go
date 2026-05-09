package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"gofoundry/db/mongodb"

	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CacheEntity struct {
	Key      string    `bson:"_id"`
	Value    []byte    `bson:"data"`
	ExpireAt time.Time `bson:"_expireAt"`
}

// must: while use mongo as cache or do manually
func MongoIndex() mongo.IndexModel {
	return mongo.IndexModel{
		Keys:    bson.D{{Key: "_expireAt", Value: 1}},
		Options: options.Index().SetName("_expireAt_").SetExpireAfterSeconds(0),
	}
}

type ICache interface {
	Set(ctx *context.Context, key string, value []byte, ttlInMinutes ...int) error
	GetMany(ctx *context.Context, key ...string) ([][]byte, error)
	Get(ctx *context.Context, key string) ([]byte, error)
	Delete(ctx *context.Context, keys ...string) error
	Type() string
}

type IMongoCache interface {
	Mongo() *mongo.Collection
	ICache
}

type IRedisCache interface {
	Redis() *redis.Client
	ICache
}

type Client struct {
	_mongo *mongo.Collection
	_redis *redis.Client
	ttl    *int //in minutes
}

// Redis implements IRedisCache.
func (c *Client) Redis() *redis.Client {
	return c._redis
}

// Mongo implements IMongoCache.
func (c *Client) Mongo() *mongo.Collection {
	return c._mongo
}

// Type implements ICacheClient.
func (c *Client) Type() string {
	if c._redis != nil && c._mongo != nil {
		return "Redis: Primary | Mongo: Secondary"
	} else if c._redis != nil {
		return "Redis"
	} else if c._mongo != nil {
		return "Mongo"
	}
	return "Nothing"
}

// Delete implements ICacheClient.
func (c *Client) Delete(ctx *context.Context, keys ...string) error {
	if len(keys) < 1 {
		return nil
	}
	chanErr1 := make(chan error)
	chanErr2 := make(chan error)
	go func() {
		chanErr1 <- func() error {
			if c._redis != nil {
				return c._redis.Del(*ctx, keys...).Err()
			}
			return nil
		}()
		close(chanErr1)
	}()
	go func() {
		chanErr2 <- func() error {
			if c._mongo != nil { //fallback
				_, err := c._mongo.DeleteMany(*ctx, bson.M{"_id": bson.M{"$in": keys}})
				return err
			}
			return nil
		}()
		close(chanErr2)
	}()
	err := errors.Join(<-chanErr1, <-chanErr2)
	return err
}

// Get implements ICacheClient.
func (c *Client) Get(ctx *context.Context, key string) (b []byte, e error) {
	if key == "" {
		return nil, errors.New("invalid key")
	} else if c._redis != nil {
		b, e = c._redis.Get(*ctx, key).Bytes()
	}
	if c._mongo != nil && len(b) < 1 { //fallback
		var data CacheEntity
		if err := c._mongo.FindOne(*ctx, bson.M{"_id": key}).Decode(&data); err != nil {
			return nil, errors.Join(e, err)
		}
		return data.Value, nil
	}
	return
}

// Get implements ICacheClient.
func (c *Client) GetMany(ctx *context.Context, keys ...string) (b [][]byte, e error) {
	if len(keys) < 1 {
		return nil, errors.New("empty key list")
	} else if c._redis != nil {
		if v, err := c._redis.MGet(*ctx, keys...).Result(); err != nil {
			e = err //fallback
		} else if len(v) > 0 {
			b = make([][]byte, len(v))
			for i := range v {
				if v[i] == nil {
					b[i] = nil
					continue
				}
				switch val := v[i].(type) {
				case []byte:
					b[i] = val
				}
			}
		}
	}
	if c._mongo != nil && len(b) < 1 { //fallback
		if data, err := mongodb.ExecuteTCursor[CacheEntity](ctx, func() (*mongo.Cursor, error) {
			return c._mongo.Find(*ctx, bson.M{"_id": bson.M{"$in": keys}})
		}); err != nil {
			return nil, errors.Join(e, err)
		} else if len(data) > 0 {
			b = make([][]byte, len(data))
			for i := range data {
				b[i] = data[i].Value
			}
		}
	}
	return
}

// Set implements ICacheClient.
func (c *Client) Set(ctx *context.Context, key string, value []byte, ttlInMinutes ...int) error {
	if key == "" {
		return errors.New("invalid key")
	} else if len(ttlInMinutes) < 1 && c.ttl != nil {
		ttlInMinutes = append(ttlInMinutes, *c.ttl)
	} else if len(ttlInMinutes) < 1 {
		ttlInMinutes = append(ttlInMinutes, 20) //20mins
	}

	chanErr1 := make(chan error)
	chanErr2 := make(chan error)
	go func() {
		chanErr1 <- func() error {
			if c._redis != nil {
				return c._redis.Set(*ctx, key, value, time.Duration(ttlInMinutes[0])*time.Minute).Err()
			}
			return nil
		}()
		close(chanErr1)
	}()
	go func() {
		chanErr2 <- func() error {
			if c._mongo != nil { //fallback
				res, err := c._mongo.UpdateByID(*ctx, key, bson.M{"$set": CacheEntity{
					Key:      key,
					Value:    value,
					ExpireAt: time.Now().Add(time.Duration(ttlInMinutes[0]) * time.Minute),
				}}, options.Update().SetUpsert(true))
				if err != nil {
					return err
				} else if res == nil || (res.ModifiedCount < 1 && res.UpsertedCount < 1) {
					return errors.New("uncached")
				}
			}
			return nil
		}()
		close(chanErr2)
	}()
	err := errors.Join(<-chanErr1, <-chanErr2)
	return err
}

var _default ICache

func Default() ICache {
	return _default
}

func GetClient() *Client {
	if v, ok := _default.(*Client); ok && v != nil {
		return v
	}
	return nil
}

func NewMongoCache(collectionName string, db *mongo.Database, ttlInMinutes ...int) IMongoCache {
	var ttl *int
	if len(ttlInMinutes) > 0 {
		ttl = &ttlInMinutes[0]
	}
	mc := &Client{
		_mongo: db.Collection(collectionName),
		ttl:    ttl,
	}
	_default = mc
	return mc
}

func (c *Client) AddMongo(collectionName string, db *mongo.Database) {
	c._mongo = db.Collection(collectionName)
}

func NewRedisCache(ttlInMinutes ...int) IRedisCache {
	var ttl *int
	if len(ttlInMinutes) > 0 {
		ttl = &ttlInMinutes[0]
	}
	rc := &Client{
		_redis: redis.NewClient(&redis.Options{}),
		ttl:    ttl,
	}
	_default = rc
	return rc
}

func (c *Client) AddRedis() {
	//TODO
}

func NewCache(client func() ICache) ICache {
	_default = client()
	return _default
}

type CacheOptions struct {
	Client ICache
	TTL    int //in minutes
}

func Get[T any](ctx *context.Context, key string, client ...ICache) (*T, error) {
	var data T
	cc := _default
	if len(client) > 0 && client[0] != nil {
		cc = client[0]
	} else if cc == nil {
		return nil, errors.New("cache unconfigured")
	}
	if bytes, err := cc.Get(ctx, key); err != nil {
		return nil, err
	} else if len(bytes) < 1 {
		return nil, errors.New("not found")
	} else if err := msgpack.Unmarshal(bytes, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func GetMany[T any](ctx *context.Context, keys []string, client ...ICache) ([]T, error) {
	cc := _default
	if len(client) > 0 && client[0] != nil {
		cc = client[0]
	} else if cc == nil {
		return nil, errors.New("cache unconfigured")
	}
	if bytes, err := cc.GetMany(ctx, keys...); err != nil {
		return nil, err
	} else if len(bytes) > 0 {
		result := make([]T, len(bytes))
		for i := range bytes {
			msgpack.Unmarshal(bytes[i], &result[i])
		}
		return result, nil
	}
	return nil, errors.New("empty result")
}

func Set[T any](ctx *context.Context, key string, value *T, option ...CacheOptions) error {
	var ttl []int
	client := _default
	if len(option) > 0 && option[0].Client != nil {
		client = option[0].Client
	} else if client == nil {
		return errors.New("cache unconfigured")
	}
	if len(option) > 0 && option[0].TTL > 0 {
		ttl = []int{option[0].TTL}
	}
	bytes, err := msgpack.Marshal(value)
	if err != nil {
		return err
	} else if err = client.Set(ctx, key, bytes, ttl...); err != nil {
		return err
	}
	return nil
}

func SetJson(ctx *context.Context, key string, value any, option ...CacheOptions) error {
	var ttl []int
	client := _default
	if len(option) > 0 && option[0].Client != nil {
		client = option[0].Client
	} else if client == nil {
		return errors.New("cache unconfigured")
	}
	if len(option) > 0 && option[0].TTL > 0 {
		ttl = []int{option[0].TTL}
	}
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	} else if err = client.Set(ctx, key, bytes, ttl...); err != nil {
		return err
	}
	return nil
}

func Delete(ctx *context.Context, keys []string, client ...ICache) error {
	cc := _default
	if len(client) > 0 && client[0] != nil {
		cc = client[0]
	} else if cc == nil {
		return errors.New("cache unconfigured")
	}
	return cc.Delete(ctx, keys...)
}

type MongoOptions struct {
	Client IMongoCache
	TTL    int //in minutes
}

func SetBson(ctx *context.Context, key string, updates []bson.M, option ...MongoOptions) error {
	var ttl []int
	var cache IMongoCache
	if key == "" || len(updates) < 1 {
		return errors.New("invalid key or empty updates")
	} else if len(option) > 0 && option[0].Client != nil {
		cache = option[0].Client
	} else if c := GetClient(); c != nil {
		cache = c
		if c.ttl != nil {
			ttl = []int{*c.ttl}
		}
	} else {
		return errors.New("cache unconfigured")
	}
	if len(option) > 0 && option[0].TTL > 0 {
		ttl = []int{option[0].TTL}
	} else if len(ttl) < 1 {
		ttl = []int{20} //20mins
	}
	if mc := cache.Mongo(); mc != nil {
		wm := []mongo.WriteModel{}
		for i := range updates {
			wm = append(wm, mongo.NewUpdateOneModel().SetFilter(bson.M{"_id": key}).SetUpdate(updates[i]).SetUpsert(true))
		}
		wm = append(wm, mongo.NewUpdateOneModel().SetFilter(bson.M{"_id": key}).SetUpdate(bson.M{"$set": bson.M{"_expireAt": time.Now().Add(time.Duration(ttl[0]) * time.Minute)}}))
		res, err := mc.BulkWrite(*ctx, wm, options.BulkWrite().SetOrdered(true))
		if err != nil {
			return err
		} else if res == nil || (res.ModifiedCount < 1 && res.UpsertedCount < 1) {
			return errors.New("uncached")
		}
		return nil
	}
	return errors.New("cache (mongo) unconfigured")
}

func GetBson[T any](ctx *context.Context, key string, stages []bson.M, option ...MongoOptions) ([]T, error) {
	var cache IMongoCache
	if key == "" {
		return nil, errors.New("invalid key")
	} else if len(option) > 0 && option[0].Client != nil {
		cache = option[0].Client
	} else if c := GetClient(); c != nil {
		cache = c
	} else {
		return nil, errors.New("cache unconfigured")
	}
	if mc := cache.Mongo(); mc != nil {
		return mongodb.ExecuteTCursor[T](ctx, func() (*mongo.Cursor, error) {
			return mc.Aggregate(*ctx, append([]bson.M{{"$match": bson.M{"_id": key}}}, stages...), options.Aggregate().SetCollation(&options.Collation{
				Locale:   "en",
				Strength: 2,
			}))
		})
	}
	return nil, errors.New("cache (mongo) unconfigured")
}
