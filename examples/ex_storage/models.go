package ex_storage

import (
	"github.com/simpmeth/geojsonrm"
)

type Address struct {
	ID       uint `gorm:"primaryKey"`
	Address  string
	GeoPoint geojsonrm.Point
}

type Zone struct {
	ID         uint `gorm:"primaryKey"`
	Title      string
	GeoPolygon geojsonrm.Polygon
}

type Route struct {
	ID       uint `gorm:"primaryKey"`
	Title    string
	GeoRoute geojsonrm.LineString
}
