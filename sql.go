package geojsonrm

import (
	"bytes"
	"database/sql/driver"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"github.com/twpayne/go-geom/encoding/geojson"
	"github.com/twpayne/go-geom/encoding/wkt"
	"strconv"
)

var (
	ErrUnexpectedGeometryType = errors.New("unexpected geometry type")
	ErrUnexpectedValueType    = errors.New("unexpected value type")
)

var SRID = 4326

type (
	Geometry[T geom.T] struct{ Geom geom.T }

	Point              = Geometry[*geom.Point]
	LineString         = Geometry[*geom.LineString]
	Polygon            = Geometry[*geom.Polygon]
	MultiPoint         = Geometry[*geom.MultiPoint]
	MultiLineString    = Geometry[*geom.MultiLineString]
	MultiPolygon       = Geometry[*geom.MultiPolygon]
	GeometryCollection = Geometry[*geom.GeometryCollection]
)

func New[T geom.T](geom T) Geometry[T] { return Geometry[T]{geom} }

func (g *Geometry[T]) Scan(value interface{}) (err error) {
	var (
		wkb []byte
		ok  bool
	)
	switch v := value.(type) {
	case string:
		wkb, err = hex.DecodeString(v)
	case []byte:
		wkb = v
	default:
		return ErrUnexpectedGeometryType
	}
	if err != nil {
		return err
	}
	geometryT, err := ewkb.Unmarshal(wkb)
	if err != nil {
		return err
	}
	g.Geom, ok = geometryT.(T)
	if !ok {
		return ErrUnexpectedValueType
	}

	return
}

func (g Geometry[T]) Value() (driver.Value, error) {
	if g.Geom == nil {
		return nil, nil
	}

	sb := &bytes.Buffer{}
	if err := ewkb.Write(sb, binary.LittleEndian, g.Geom); err != nil {
		return nil, err
	}
	return hex.EncodeToString(sb.Bytes()), nil
}

func (g Geometry[T]) GormDataType() string {
	srid := strconv.Itoa(SRID)

	switch any(g).(type) {
	case Geometry[*geom.Point]:
		return "Geometry(Point, " + srid + ")"
	case Geometry[*geom.LineString]:
		return "Geometry(LineString, " + srid + ")"
	case Geometry[*geom.Polygon]:
		return "Geometry(Polygon, " + srid + ")"
	case Geometry[*geom.MultiPoint]:
		return "Geometry(MultiPoint, " + srid + ")"
	case Geometry[*geom.MultiLineString]:
		return "Geometry(MultiLineString, " + srid + ")"
	case Geometry[*geom.MultiPolygon]:
		return "Geometry(MultiPolygon, " + srid + ")"
	case Geometry[*geom.GeometryCollection]:
		return "Geometry(GeometryCollection, " + srid + ")"
	default:
		return "geometry"
	}
}

func (g Geometry[T]) String() string {
	if geomWkt, err := wkt.Marshal(g.Geom); err == nil {
		return geomWkt
	}
	return fmt.Sprintf("cannot marshal geometry: %T", g.Geom)
}

func (g *Geometry[T]) MarshalJSON() ([]byte, error) {
	return geojson.Marshal(g.Geom)
}

func (g *Geometry[T]) UnmarshalJSON(data []byte) error {
	return geojson.Unmarshal(data, &g.Geom)
}
