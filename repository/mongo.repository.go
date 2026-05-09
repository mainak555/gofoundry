package repository

import (
	"context"
	"reflect"
	"time"

	libmongo "mongodb"
	"repository/interfaces"
	"util"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	_collection  *mongo.Collection
	_constructor func() any
}

type MongoTRepository[T interface{}] struct {
	_base interfaces.IMongoRepository
}

func NewMongoTRepositoryOfCollection[T interface{}](collectionName string, mClient libmongo.IMongoClient) interfaces.TMongoRepository[T] {
	return &MongoTRepository[T]{
		_base: NewMongoRepository[T](collectionName, mClient),
	}
}

func NewMongoTRepository[T interface{}](mClient libmongo.IMongoClient) interfaces.TMongoRepository[T] {
	var entity T
	collectionName, err := mClient.GetCollectionName(entity)
	if err != nil {
		util.PanicError(err)
	}
	return &MongoTRepository[T]{
		_base: NewMongoRepository[T](*collectionName, mClient),
	}
}

func NewMongoRepository[T interface{}](collectionName string, mClient libmongo.IMongoClient) interfaces.IMongoRepository {
	return &MongoRepository{
		_collection: mClient.GetDatabase().Collection(collectionName),
		_constructor: func() any {
			var entity T
			return &entity
		},
	}
}

func (r *MongoTRepository[T]) Base() interfaces.IMongoRepository {
	return r._base
}

func (r *MongoTRepository[T]) Collection() *mongo.Collection {
	return r._base.Collection()
}

func (r *MongoRepository) Collection() *mongo.Collection {
	return r._collection
}

/*==========================================================GET==================================================*/

func (r *MongoTRepository[T]) Get(ctx *context.Context, query map[string]any) ([]T, error) {
	return r.GetByPage(ctx, query, 1, 100)
}

func (r *MongoTRepository[T]) GetByPage(ctx *context.Context, query map[string]any, pageNo, pageSize int64) ([]T, error) {
	return r.GetWithOptions(ctx, query, MongoPaginationOptions(pageNo, pageSize))
}

func (r *MongoTRepository[T]) GetWithOptions(ctx *context.Context, query bson.M, mongoOptions *options.FindOptions) ([]T, error) {
	var _result []T
	err := _mongoGetWithOptions(ctx, query, r.Collection(), func(decode func(val interface{}) error) error {
		var entry T
		if err := decode(&entry); err != nil {
			return err
		}
		_result = append(_result, entry)
		return nil
	}, mongoOptions)
	return _result, err
}

/*---------------------------------------------------------------------------------------------------------------------------------------*/

func (r *MongoRepository) Get(ctx *context.Context, query map[string]any) ([]any, error) {
	return r.GetByPage(ctx, query, 1, 100)
}

func (r *MongoRepository) GetByPage(ctx *context.Context, query map[string]any, pageNo, pageSize int64) ([]any, error) {
	return r.GetWithOptions(ctx, query, MongoPaginationOptions(pageNo, pageSize))
}

func (r *MongoRepository) GetWithOptions(ctx *context.Context, query bson.M, mongoOptions *options.FindOptions) ([]any, error) {
	var _result []any
	err := _mongoGetWithOptions(ctx, query, r._collection, func(decode func(val interface{}) error) error {
		entry := r._constructor()
		if err := decode(entry); err != nil {
			return err
		}
		_result = append(_result, entry)
		return nil
	}, mongoOptions)
	return _result, err
}

func _mongoGetWithOptions(ctx *context.Context, query bson.M, collection *mongo.Collection,
	fn func(func(val interface{}) error) error, mongoOptions *options.FindOptions) error {

	if len(query) > 0 {
		query = bson.M{"$and": []bson.M{{"_deleted": false}, query}}
	} else {
		query["_deleted"] = false
	}

	if mongoOptions == nil {
		mongoOptions = &options.FindOptions{}
	}
	//mongoOptions.SetMaxTime(time.Minute * 1)

	cur, err := collection.Find(*ctx, query, mongoOptions)
	if err != nil {
		return err
	}

	defer cur.Close(*ctx)
	for cur.Next(*ctx) {
		if err := fn(cur.Decode); err != nil {
			return err
		}
	}
	return nil
}

func MongoPaginationOptions(pageNo, pageSize int64) *options.FindOptions {
	skip := pageNo*pageSize - pageSize
	options := &options.FindOptions{
		Limit: &pageSize,
		Skip:  &skip,
	}
	return options
}

/*==========================================================GET ONE==================================================*/

func (r *MongoTRepository[T]) GetOne(ctx *context.Context, id string) (*T, error) {
	res := _getOne(ctx, id, r.Collection())
	var dbo T
	err := res.Decode(&dbo)
	return &dbo, err
}

func (r *MongoRepository) GetOne(ctx *context.Context, id string) (interface{}, error) {
	res := _getOne(ctx, id, r._collection)
	dbo := r._constructor()
	err := res.Decode(dbo) //& not required as ctor() => &T
	return dbo, err
}

func _getOne(ctx *context.Context, id string, collection *mongo.Collection) *mongo.SingleResult {
	_id, _ := primitive.ObjectIDFromHex(id)
	return collection.FindOne(*ctx, bson.M{"_id": _id, "_deleted": false})
}

/*=======================================================CREATE ONE==================================================*/

func (r *MongoTRepository[T]) CreateOne(ctx *context.Context, obj *T) error {
	return r._base.CreateOne(ctx, obj)
}

func (r *MongoRepository) CreateOne(ctx *context.Context, obj interface{}) error {
	/*_id := reflect.ValueOf(obj).Elem().FieldByName("Id").Interface().(primitive.ObjectID)
	//reflect.ValueOf(obj).Elem().FieldByName("CreatedAt").Set(reflect.ValueOf(time.Now().UTC()))*/
	*reflect.ValueOf(obj).Elem().FieldByName("CreatedAt").Addr().Interface().(*time.Time) = time.Now().UTC()

	result, err := r._collection.InsertOne(*ctx, obj)
	if err != nil {
		return err
	}

	_id := result.InsertedID.(primitive.ObjectID)
	//r.Collection.UpdateByID(*ctx, bson.M{"_id": _id}, addVirtualDelete)
	reflect.ValueOf(obj).Elem().FieldByName("Id").Set(reflect.ValueOf(_id))
	return nil
}

/*========================================================UPDATE ONE==================================================*/

func (r *MongoTRepository[T]) UpdateOne(ctx *context.Context, id string, obj *T) (int64, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return -1, err
	}
	reflect.ValueOf(obj).Elem().FieldByName("Id").Set(reflect.ValueOf(_id))
	/* NOT REQUIRED : as taken care at base
	time := time.Now().UTC()
	reflect.ValueOf(obj).Elem().FieldByName("UpdatedAt").Set(reflect.ValueOf(&time))
	*/
	var doc bson.D
	if data, err := bson.Marshal(*obj); err != nil {
		return -1, err
	} else if err := bson.Unmarshal(data, &doc); err != nil {
		return -1, err
	}
	return r._base.UpdateOne(ctx, id, []bson.M{
		{"$set": doc},
	})
}

func (r *MongoRepository) UpdateOne(ctx *context.Context, id string, bsonObj []primitive.M) (int64, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return -1, err
	}

	bsonObj = append(bsonObj, bson.M{
		"$set": bson.M{
			"_updatedAt": time.Now().UTC(),
		},
	})

	filter := bson.M{"$and": []bson.M{{"_id": _id}, {"_deleted": false}}} //Allow Update Only undeleted marked records
	result, err := r._collection.UpdateOne(*ctx, filter, bsonObj)
	if result != nil {
		return result.ModifiedCount, err
	}
	return -1, err
}

/*========================================================UPDATE MANY=================================================*/

func (r *MongoTRepository[T]) UpdateMany(ctx *context.Context, query bson.M, bsonObj []bson.M) (int64, error) {
	return r._base.UpdateMany(ctx, query, bsonObj)
}

func (r *MongoRepository) UpdateMany(ctx *context.Context, query bson.M, bsonObj []bson.M) (int64, error) {
	bsonObj = append(bsonObj, bson.M{
		"$set": bson.M{
			"_updatedAt": time.Now().UTC(),
		},
	})

	filter := bson.M{
		"$and": []bson.M{
			{"_deleted": false},
			query,
		},
	} //Allow Update Only undeleted marked records
	result, err := r._collection.UpdateMany(*ctx, filter, bsonObj)
	if result != nil {
		return result.ModifiedCount, err
	}
	return -1, err
}

/*====================================================SOFT DELETE==================================================*/

func (r *MongoTRepository[T]) DeleteOne(ctx *context.Context, id string, params func() map[string]any) (int64, error) {
	return r._base.DeleteOne(ctx, id, params)
}

func (r *MongoRepository) DeleteOne(ctx *context.Context, id string, params func() map[string]any) (int64, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return -1, err
	}
	filter := bson.M{"_id": _id}
	return r.Delete(ctx, filter, params)
}

func (r *MongoTRepository[T]) Delete(ctx *context.Context, filter map[string]any, params func() map[string]any) (int64, error) {
	return r._base.Delete(ctx, filter, params)
}

func (r *MongoRepository) Delete(ctx *context.Context, filter map[string]any, params func() map[string]any) (int64, error) {
	update := bson.M{
		"$set": bson.M{
			"_deleted":   true,
			"_deletedAt": time.Now().UTC(),
		},
	}

	query := bson.M{"$and": []bson.M{
		{"_deleted": false}, filter},
	} //Allow Update Only undeleted marked records

	if params != nil {
		fields := params()
		for k, v := range fields {
			update["$set"].(bson.M)[k] = v
		}
	}

	result, err := r._collection.UpdateMany(*ctx, query, update)
	if err != nil {
		return -1, err
	}
	return result.ModifiedCount, nil
}

/*======================================================HARD DELETE===============================================*/

func (r *MongoTRepository[T]) EraseOne(ctx *context.Context, id string) (int64, error) {
	return r._base.EraseOne(ctx, id)
}

func (r *MongoRepository) EraseOne(ctx *context.Context, id string) (int64, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return -1, err
	}
	filter := bson.M{"_id": _id}
	return r.Erase(ctx, filter)
}

func (r *MongoTRepository[T]) Erase(ctx *context.Context, filter map[string]any) (int64, error) {
	return r._base.Erase(ctx, filter)
}

func (r *MongoRepository) Erase(ctx *context.Context, filter map[string]any) (int64, error) {
	result, err := r._collection.DeleteMany(*ctx, filter)
	if err != nil {
		return -1, err
	}
	return result.DeletedCount, nil
}
