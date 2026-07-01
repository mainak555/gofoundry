package mongodb

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BaseIdEntity struct {
	Id primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
}

type BaseEntity struct {
	BaseIdEntity `bson:",inline"`
	CreatedAt    time.Time  `bson:"_createdAt" json:"createdAt"`
	UpdatedAt    *time.Time `bson:"_updatedAt,omitempty" json:"updatedAt,omitempty"`
	DeletedAt    *time.Time `bson:"_deletedAt,omitempty" json:"deletedAt,omitempty"`
	CreatedBy    string     `bson:"_createdBy,omitempty" json:"createdBy,omitempty"`
	UpdatedBy    string     `bson:"_updatedBy,omitempty" json:"updatedBy,omitempty"`
	DeletedBy    string     `bson:"_deletedBy,omitempty" json:"deletedBy,omitempty"`
	Deleted      bool       `bson:"_deleted" json:"deleted"`
}
