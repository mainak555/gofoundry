package mongodb

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"util"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"golang.org/x/exp/slices"
)

type IMongoClient interface {
	Ping() error
	Disconnect() error
	GetClient() *mongo.Client
	GetDatabase() *mongo.Database
	IsCollectionExist(string) bool
	Connect(*context.Context) (IMongoClient, error)
	GetCollectionName(interface{}) (*string, error)
	ConfigureCollections(*context.Context, ...func() (string, interface{}, []mongo.IndexModel)) error
}

type MongoClient struct {
	_dbName               string
	_client               *mongo.Client
	_database             *mongo.Database
	_registeredColletions map[string]string
	_clientOpts           *options.ClientOptions
	_dbOpts               *options.DatabaseOptions
}

func NewMongoClient(ctx *context.Context, conStr string, dbName string,
	_clientOpts func(*options.ClientOptions) *options.ClientOptions, _dbOps func(*options.DatabaseOptions) *options.DatabaseOptions) (IMongoClient, error) {

	mc := &MongoClient{
		_dbOpts:     options.Database().SetReadPreference(readpref.Secondary()),
		_clientOpts: options.Client().ApplyURI(conStr).SetWriteConcern(writeconcern.Majority()),
	}

	if _clientOpts != nil {
		mc._clientOpts = _clientOpts(mc._clientOpts)
	}

	if _dbOps != nil {
		mc._dbOpts = _dbOps(mc._dbOpts)
	}

	mc._dbName = dbName
	return mc.Connect(ctx)
}

func (c *MongoClient) Connect(ctx *context.Context) (IMongoClient, error) {
	client, err := mongo.Connect(*ctx, c._clientOpts)
	if err != nil {
		return nil, err
	}
	c._client = client
	c._database = client.Database(c._dbName, c._dbOpts)
	return c, nil
}

func (c *MongoClient) Ping() error {
	return c._client.Ping(context.TODO(), nil)
}

func (c *MongoClient) Disconnect() error {
	return c._client.Disconnect(context.TODO())
}

func (c *MongoClient) IsCollectionExist(key string) bool {
	_, ok := c._registeredColletions[key]
	return ok
}

func (c *MongoClient) GetCollectionName(key interface{}) (*string, error) {
	var collName string
	if reflect.TypeOf(key).Kind() == reflect.Struct {
		entity := reflect.TypeOf(key).Name()
		collName = c._registeredColletions[entity]
		if collName == "" {
			return nil, fmt.Errorf("mongo collection for entity [%v] not registered", entity)
		}
	} else if reflect.TypeOf(key).Kind() == reflect.String {
		if v, ok := c._registeredColletions[key.(string)]; ok {
			collName = v
		} else {
			return nil, fmt.Errorf("mongo collection [%v] not registered", key)
		}
	} else {
		return nil, errors.New("key should be type of struct/string")
	}
	return &collName, nil
}

func (c *MongoClient) GetClient() *mongo.Client {
	return c._client
}

func (c *MongoClient) GetDatabase() *mongo.Database {
	return c._database
}

func (c *MongoClient) ConfigureCollections(ctx *context.Context, collections ...func() (string, interface{}, []mongo.IndexModel)) error {
	c._registeredColletions = make(map[string]string)
	collNames, err := c._database.ListCollectionNames(*ctx, bson.M{})
	if err != nil {
		return err
	}

	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
	for _, cf := range collections {
		collName, entity, indxModels := cf()

		if entity != nil && reflect.TypeOf(entity).Kind() != reflect.Struct {
			util.PanicError(errors.New("entity is not type of struct"))
		}

		if !slices.Contains(collNames, collName) {
			c._database.CreateCollection(*ctx, collName)
			collection := c._database.Collection(collName)
			if indxModels != nil || len(indxModels) > 0 {
				_, indexErr := collection.Indexes().CreateMany(context.TODO(), indxModels, opts)
				if indexErr != nil {
					return indexErr
				}
				fmt.Println("Index Created for:", collName)
			}
		} else {
			fmt.Println("Verifying Indexes for:", collName)
			collection := c._database.Collection(collName)
			indexView := collection.Indexes()
			lstOpts := options.ListIndexes().SetMaxTime(10 * time.Second)

			var result []bson.M
			if cur, err := indexView.List(*ctx, lstOpts); err != nil {
				return err
			} else {
				if err := cur.All(*ctx, &result); err != nil {
					return err
				}
			}

			var indexes []string
			for _, item := range result {
				indexes = append(indexes, item["name"].(string))
			}

			for _, im := range indxModels {
				if slices.Contains(indexes, *im.Options.Name) {
					continue
				}
				//Create Newly Added Index
				_, indexErr := collection.Indexes().CreateOne(context.TODO(), im, opts)
				if indexErr != nil {
					return indexErr
				}
				fmt.Println("- Index Added:", *im.Options.Name)
			}
		}

		if entity != nil && reflect.TypeOf(entity).Kind() == reflect.Struct {
			entityName := reflect.TypeOf(entity).Name()
			c._registeredColletions[entityName] = collName
			continue
		}
		c._registeredColletions[collName] = collName
	}
	return nil
}
