package mongodb

import (
	"errors"

	"github.com/mainak555/gofoundry/db"
)

type GeoSpatialType string

const (
	GEOSPATIAL_POINT      GeoSpatialType = "Point"
	GEOSPATIAL_POLYGON    GeoSpatialType = "Polygon"
	GEOSPATIAL_LINESTRING GeoSpatialType = "LineString"

	GEOSPATIAL_INDEX_2D       string = "2d"
	GEOSPATIAL_INDEX_2DSPHERE string = "2dsphere"
)

func (e GeoSpatialType) string() string {
	switch e {
	case GEOSPATIAL_POINT:
		return "Point"
	case GEOSPATIAL_POLYGON:
		return "Polygon"
	case GEOSPATIAL_LINESTRING:
		return "LineString"
	}
	panic(errors.New("Invalid"))
}

type BaseGeoSpatialEntity struct {
	Type GeoSpatialType `bson:"type"`
}

type GeoSpatialPointEntity struct {
	BaseGeoSpatialEntity `bson:",inline"`
	Coordinates          []float64 `bson:"coordinates" json:"coordinates"`
}

func NewGeoSpatialPointEntity(latitude, longitude float64) *GeoSpatialPointEntity {
	return &GeoSpatialPointEntity{
		BaseGeoSpatialEntity: BaseGeoSpatialEntity{
			Type: GEOSPATIAL_POINT,
		},
		Coordinates: []float64{
			longitude, latitude,
		},
	}
}

func NewMongoGeoLocationEntity(geoLocation *db.GeoCoordinateEntity) *GeoSpatialPointEntity {
	if geoLocation != nil {
		return NewGeoSpatialPointEntity(geoLocation.Latitude, geoLocation.Longitude)
	}
	return nil
}

func (loc *GeoSpatialPointEntity) GetLatLong() *db.GeoCoordinateEntity {
	if loc != nil && loc.Coordinates != nil && len(loc.Coordinates) == 2 {
		return &db.GeoCoordinateEntity{
			Latitude:  loc.Coordinates[1],
			Longitude: loc.Coordinates[0],
		}
	}
	return nil
}

type GeoSpatialLineStringEntity struct {
	BaseGeoSpatialEntity `bson:",inline"`
	Coordinates          [][]float64 `bson:"coordinates" json:"coordinates"`
}

type GeoSpatialPolygonEntity struct {
	BaseGeoSpatialEntity `bson:",inline"`
	Coordinates          [][][]float64 `bson:"coordinates" json:"coordinates"`
}
