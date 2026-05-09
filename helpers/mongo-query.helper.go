package helpers

import (
	"gofoundry/util"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MongoQueryDeletedBy(deletedBy string) func() map[string]any {
	return func() map[string]any {
		return map[string]any{
			"deletedBy": deletedBy,
		}
	}
}

func MongoQueryInByIds(ids ...string) bson.M {
	Ids, _ := util.StringArrayToObjectIDArray(ids...)
	return bson.M{
		"_id": bson.M{
			"$in": Ids,
		},
	}
}

func GetNonNullObjectID(input *string) primitive.ObjectID {
	if input != nil && *input != "" {
		if val, err := util.StringArrayToObjectIDArray(*input); err == nil {
			return val[0]
		}
	}
	return primitive.NilObjectID
}

func GetNullableObjectID(input *string) *primitive.ObjectID {
	if input != nil && *input != "" {
		if val, err := util.StringArrayToObjectIDArray(*input); err == nil {
			return &val[0]
		}
	}
	return nil
}
