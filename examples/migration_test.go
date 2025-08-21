package examples

import (
	"testing"

	"github.com/simpmeth/geojsonrm"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type TempTableWithGeometry[Geometry any] struct {
	gorm.Model
	Geometry Geometry
}

func (temp TempTableWithGeometry[Geometry]) TableName() string {
	return "temp_table_with_geometry"
}

func TestMigrate(t *testing.T) {
	tests := []struct {
		model              any
		expectGeometryType string
	}{
		{
			model:              TempTableWithGeometry[geojsonrm.Point]{},
			expectGeometryType: "geometry(Point,4326)",
		},
		{
			model:              TempTableWithGeometry[geojsonrm.LineString]{},
			expectGeometryType: "geometry(LineString,4326)",
		},
		{
			model:              TempTableWithGeometry[geojsonrm.Polygon]{},
			expectGeometryType: "geometry(Polygon,4326)",
		},
		{
			model:              TempTableWithGeometry[geojsonrm.MultiPoint]{},
			expectGeometryType: "geometry(MultiPoint,4326)",
		},
		{
			model:              TempTableWithGeometry[geojsonrm.MultiLineString]{},
			expectGeometryType: "geometry(MultiLineString,4326)",
		},
		{
			model:              TempTableWithGeometry[geojsonrm.MultiPolygon]{},
			expectGeometryType: "geometry(MultiPolygon,4326)",
		},
		{
			model:              TempTableWithGeometry[geojsonrm.GeometryCollection]{},
			expectGeometryType: "geometry(GeometryCollection,4326)",
		},
	}

	for _, test := range tests {
		t.Run("migrate table with "+test.expectGeometryType, func(t *testing.T) {
			migrator := db.Migrator()

			err := migrator.CreateTable(test.model)
			require.NoError(t, err)

			defer func() {
				_ = migrator.DropTable(test.model)
			}()

			// check if column geometry exists
			require.True(t, migrator.HasColumn(test.model, "geometry"))

			// find column geometry
			columns, err := migrator.ColumnTypes(test.model)
			require.NoError(t, err)

			for _, column := range columns {
				if column.Name() != "geometry" {
					continue
				}

				columnType, ok := column.ColumnType()
				require.True(t, ok)

				require.Equal(t, test.expectGeometryType, columnType)

				return
			}
		})
	}
}
